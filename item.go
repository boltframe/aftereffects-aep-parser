package aep

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/rioam2/rifx"
)

// ItemTypeName denotes the type of item. See: http://docs.aenhancers.com/items/item/#item-typename
type ItemTypeName string

const (
	// ItemTypeFolder denotes a Folder item which may contain additional items
	ItemTypeFolder ItemTypeName = "Folder"
	// ItemTypeComposition denotes a Composition item which has a dimension, length, framerate and child layers
	ItemTypeComposition ItemTypeName = "Composition"
	// ItemTypeFootage denotes an AVItem that has a source (eg: an image or video file)
	ItemTypeFootage ItemTypeName = "Footage"
)

// FootageType denotes the type of footage of an AVItem (eg: Solid, Placeholder, ...)
type FootageType uint16

const (
	// FootageTypeSolid denotes a Solid source
	FootageTypeSolid FootageType = 0x09
	// FootageTypePlaceholder denotes a Placeholder source
	FootageTypePlaceholder FootageType = 0x02
)

// Item is a generalized object storing information about folders, compositions, or footage
type Item struct {
	Name              string
	ID                uint32
	TypeName          ItemTypeName
	FolderContents    []*Item
	FootageDimensions [2]uint16
	FootageFramerate  float64
	FootageSeconds    float64
	FootageType       FootageType
	BackgroundColor   [3]byte
}

func parseItem(itemHead *rifx.List) (*Item, error) {
	item := &Item{}
	isRoot := itemHead.Identifier == "Fold"

	// Parse item metadata
	if isRoot {
		item.ID = 0
		item.Name = "root"
		item.TypeName = ItemTypeFolder
	} else {
		nameBlock, err := itemHead.FindByType("Utf8")
		if err != nil {
			return nil, err
		}
		item.Name = nameBlock.ToString()
		type IDTA struct {
			Type      uint16
			Unknown00 [14]byte
			ID        uint32
		}
		itemDescriptor := &IDTA{}
		idtaBlock, err := itemHead.FindByType("idta")
		if err != nil {
			return nil, err
		}
		err = idtaBlock.ToStruct(itemDescriptor)
		if err != nil {
			return nil, err
		}
		item.ID = itemDescriptor.ID
		switch itemDescriptor.Type {
		case 0x01:
			item.TypeName = ItemTypeFolder
		case 0x04:
			item.TypeName = ItemTypeComposition
		case 0x07:
			item.TypeName = ItemTypeFootage
		}
	}

	// Parse unique item type information
	switch item.TypeName {
	case ItemTypeFolder:
		childItemLists := append(itemHead.SublistFilter("Item"), itemHead.SublistMerge("Sfdr").SublistFilter("Item")...)
		for _, childItemList := range childItemLists {
			childItem, err := parseItem(childItemList)
			if err != nil {
				return nil, err
			}
			item.FolderContents = append(item.FolderContents, childItem)
		}
	case ItemTypeFootage:
		pinList, err := itemHead.SublistFind("Pin ")
		if err != nil {
			return nil, err
		}
		sspcBlock, err := pinList.FindByType("sspc")
		if err != nil {
			return nil, err
		}
		type SSPC struct {
			Unknown00         [30]byte // Offset 0B
			Width             uint32   // Offset 30B
			Height            uint32   // Offset 34B
			SecondsDividend   uint32   // Offset 38B
			SecondsDivisor    uint32   // Offset 42B
			Unknown01         [10]byte // Offset 46B
			Framerate         uint32   // Offset 56B
			FramerateDividend uint16   // Offset 60B
		}
		sspcDesc := &SSPC{}
		sspcBlock.ToStruct(sspcDesc)
		item.FootageDimensions = [2]uint16{uint16(sspcDesc.Width), uint16(sspcDesc.Height)}
		item.FootageFramerate = float64(sspcDesc.Framerate) + (float64(sspcDesc.FramerateDividend) / float64(1<<16))
		item.FootageSeconds = float64(sspcDesc.SecondsDividend) / float64(sspcDesc.SecondsDivisor)

		optiBlock, err := pinList.FindByType("opti")
		if err != nil {
			return nil, err
		}
		optiData := optiBlock.Data.([]byte)
		item.FootageType = FootageType(binary.BigEndian.Uint16(optiData[4:6]))
		switch item.FootageType {
		case FootageTypeSolid:
			item.Name = fmt.Sprintf("%s", bytes.ReplaceAll(bytes.Trim(optiData[26:255], "\x00"), []byte{0}, []byte{32}))
		case FootageTypePlaceholder:
			item.Name = fmt.Sprintf("%s", bytes.ReplaceAll(bytes.Trim(optiData[10:], "\x00"), []byte{0}, []byte{32}))
		}
	case ItemTypeComposition:
		type CDTA struct {
			Unknown00         [4]byte  // Offset 0B
			FramerateDivisor  uint32   // Offset 4B
			FramerateDividend uint32   // Offset 8B
			Unknown01         [32]byte // Offset 12B
			SecondsDividend   uint32   // Offset 40B
			SecondsDivisor    uint32   // Offset 44B
			BackgroundColor   [3]byte  // Offset 48B
			Unknown03         [85]byte // Offset 51B
			Width             uint16   // Offset 136B
			Height            uint16   // Offset 138B
			Unknown04         [12]byte // Offset 140B
			Framerate         uint16   // Offset 152B
		}
		compDesc := &CDTA{}
		cdataBlock, err := itemHead.FindByType("cdta")
		if err != nil {
			return nil, err
		}
		cdataBlock.ToStruct(compDesc)
		item.FootageDimensions = [2]uint16{compDesc.Width, compDesc.Height}
		item.FootageFramerate = float64(compDesc.FramerateDividend) / float64(compDesc.FramerateDivisor)
		item.FootageSeconds = float64(compDesc.SecondsDividend) / float64(compDesc.SecondsDivisor)
		item.BackgroundColor = compDesc.BackgroundColor
	}

	return item, nil
}
