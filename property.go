package aep

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/rioam2/rifx"
)

// PropertyTypeName enumerates the value/type of a property
type PropertyTypeName uint16

const (
	// PropertyTypeBoolean denotes a boolean checkbox property
	PropertyTypeBoolean PropertyTypeName = 0x04
	// PropertyTypeOneD denotes a one-dimensional slider property
	PropertyTypeOneD PropertyTypeName = 0x02
	// PropertyTypeTwoD denotes a two-dimensional point property
	PropertyTypeTwoD PropertyTypeName = 0x06
	// PropertyTypeThreeD denotes a three-dimensional point property
	PropertyTypeThreeD PropertyTypeName = 0x12
	// PropertyTypeColor denotes a four-dimensional color property
	PropertyTypeColor PropertyTypeName = 0x05
	// PropertyTypeAngle denotes a one-dimensional angle property
	PropertyTypeAngle PropertyTypeName = 0x03
	// PropertyTypeLayerSelect denotes a single-valued layer selection property
	PropertyTypeLayerSelect PropertyTypeName = 0x00
	// PropertyTypeSelect denotes a single-valued selection property
	PropertyTypeSelect PropertyTypeName = 0x07
	// PropertyTypeGroup denotes a collection/group property
	PropertyTypeGroup PropertyTypeName = 0x0d
	// PropertyTypeCustom denotes an unknown/custom property type (default)
	PropertyTypeCustom PropertyTypeName = 0x0f
)

// String translates a property type enumeration to string
func (p PropertyTypeName) String() string {
	switch p {
	case PropertyTypeBoolean:
		return "Boolean"
	case PropertyTypeOneD:
		return "OneD"
	case PropertyTypeTwoD:
		return "TwoD"
	case PropertyTypeThreeD:
		return "ThreeD"
	case PropertyTypeColor:
		return "Color"
	case PropertyTypeAngle:
		return "Angle"
	case PropertyTypeLayerSelect:
		return "LayerSelect"
	case PropertyTypeSelect:
		return "Select"
	case PropertyTypeGroup:
		return "Group"
	default:
		return "Custom"
	}
}

// Property describes a property object of a layer or nested property
type Property struct {
	MatchName     string
	Name          string
	Index         uint32
	PropertyType  PropertyTypeName
	Properties    []*Property
	SelectOptions []string
}

func parseProperty(propData interface{}, matchName string) (*Property, error) {
	prop := &Property{}

	// Apply some sensible default values
	prop.PropertyType = PropertyTypeCustom
	prop.MatchName = matchName
	prop.Name = matchName
	switch matchName {
	case "ADBE Effect Parade":
		prop.Name = "Effects"
	}

	// Handle different types of property data
	switch propData.(type) {
	case *rifx.List:
		propHead := propData.(*rifx.List)
		// Parse sub-properties
		prop.Properties = make([]*Property, 0)
		tdgpMap, orderedMatchNames := indexedGroupToMap(propHead)
		for idx, mn := range orderedMatchNames {
			subProp, err := parseProperty(tdgpMap[mn], mn)
			if err == nil {
				subProp.Index = uint32(idx) + 1
				prop.Properties = append(prop.Properties, subProp)
			}
		}

		// Parse effect sub-properties
		if propHead.Identifier == "sspc" {
			prop.PropertyType = PropertyTypeGroup
			fnamBlock, err := propHead.FindByType("fnam")
			if err == nil {
				prop.Name = fnamBlock.ToString()
			}
			parTList := propHead.SublistMerge("parT")
			subPropMatchNames, subPropPards := pairMatchNames(parTList)
			for idx, mn := range subPropMatchNames {
				// Skip first pard entry (describes parent)
				if idx == 0 {
					continue
				}
				subProp, err := parseProperty(subPropPards[idx], mn)
				if err == nil {
					subProp.Index = uint32(idx)
					prop.Properties = append(prop.Properties, subProp)
				}
			}
		}
	case []interface{}:
		for _, entry := range propData.([]interface{}) {
			if block, ok := entry.(*rifx.Block); ok {
				switch block.Type {
				case "pdnm":
					strContent := block.ToString()
					if prop.PropertyType == PropertyTypeSelect {
						prop.SelectOptions = strings.Split(strContent, "|")
					} else if strContent != "" {
						prop.Name = strContent
					}
				case "pard":
					blockData := block.Data.([]byte)
					prop.PropertyType = PropertyTypeName(binary.BigEndian.Uint16(blockData[14:16]))
					if prop.PropertyType == 0x0a {
						prop.PropertyType = PropertyTypeOneD
					}
					pardName := fmt.Sprintf("%s", bytes.Trim(blockData[16:48], "\x00"))
					if pardName != "" {
						prop.Name = pardName
					}
				}
			}
		}
	}

	return prop, nil
}

func pairMatchNames(head *rifx.List) ([]string, [][]interface{}) {
	matchNames := make([]string, 0)
	datum := make([][]interface{}, 0)
	if head != nil {
		groupIdx := -1
		skipToNextTDMNFlag := false
		for _, block := range head.Blocks {
			if block.Type == "tdmn" {
				matchName := fmt.Sprintf("%s", bytes.Trim(block.Data.([]byte), "\x00"))
				if matchName == "ADBE Group End" || matchName == "ADBE Effect Built In Params" {
					skipToNextTDMNFlag = true
					continue
				}
				matchNames = append(matchNames, matchName)
				skipToNextTDMNFlag = false
				groupIdx++
			} else if groupIdx >= 0 && !skipToNextTDMNFlag {
				if groupIdx >= len(datum) {
					datum = append(datum, make([]interface{}, 0))
				}
				switch block.Data.(type) {
				case *rifx.List:
					datum[groupIdx] = append(datum[groupIdx], block.Data)
				default:
					datum[groupIdx] = append(datum[groupIdx], block)
				}
			}
		}
	}
	return matchNames, datum
}

func indexedGroupToMap(tdgpHead *rifx.List) (map[string]*rifx.List, []string) {
	tdgpMap := make(map[string]*rifx.List, 0)
	matchNames, contents := pairMatchNames(tdgpHead)
	for idx, matchName := range matchNames {
		tdgpMap[matchName] = contents[idx][0].(*rifx.List)
	}
	return tdgpMap, matchNames
}
