package aep

import (
	"io"
	"os"

	"github.com/rioam2/rifx"
)

// BPC enumerates bits per channel
type BPC uint8

const (
	// BPC8 enumerates 8-bits-per-channel
	BPC8 BPC = 0x00
	// BPC16 enumerates 16-bits-per-channel
	BPC16 BPC = 0x01
	// BPC32 enumerates 32-bits-per-channel
	BPC32 BPC = 0x02
)

// Project holds information about an After Effects project file
type Project struct {
	ExpressionEngine string
	Depth            BPC
	RootFolder       *Item
}

// FromReader reads and creates a new project instance from an After Effects project file
func FromReader(reader io.Reader) (*Project, error) {
	rootList, err := rifx.FromReader(reader)
	if err != nil {
		return nil, err
	}
	project, err := parseProject(rootList)
	if err != nil {
		return nil, err
	}
	return project, nil
}

// Open opens, reads and creates a new project instance from an After Effects project file
func Open(path string) (*Project, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return FromReader(f)
}

// parseProject is an internal helper to decode AEP RIFF blocks
func parseProject(root *rifx.List) (*Project, error) {
	project := &Project{}

	// Parse expression engine
	expressionEngineList, err := root.SublistFind("ExEn")
	if err != nil {
		return nil, err
	}
	project.ExpressionEngine = expressionEngineList.Blocks[0].ToString()

	// Parse project head block
	type ProjectNhed struct {
		Unknown00 [15]uint8
		Depth     BPC
	}
	nhed := &ProjectNhed{}
	nhedBlock, err := root.Find(func(b *rifx.Block) bool { return b.Type == "nhed" })
	if err != nil {
		return nil, err
	}
	nhedBlock.ToStruct(nhed)
	project.Depth = nhed.Depth

	// Parse root projet folder
	rootFolderList, err := root.SublistFind("Fold")
	if err != nil {
		return nil, err
	}
	folder, err := parseItem(rootFolderList)
	if err != nil {
		return nil, err
	}
	project.RootFolder = folder

	return project, nil
}
