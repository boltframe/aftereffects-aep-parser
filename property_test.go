package aep

import (
	"testing"
)

func TestPropertyParsing(t *testing.T) {
	project, err := Open("data/Property-01.aep")
	if err != nil {
		t.Fatal(err)
	}

	comp01 := project.RootFolder.FolderContents[0]
	textLayer := comp01.CompositionLayers[1]
	expect(t, len(textLayer.Effects), 0)
	expect(t, textLayer.Text != nil)

	expressionControlsLayer := comp01.CompositionLayers[0]
	expect(t, len(expressionControlsLayer.Effects), 7)
	expect(t, expressionControlsLayer.Text == nil)

	checkboxEffect := expressionControlsLayer.Effects[0]
	expect(t, checkboxEffect.Index, uint32(1))
	expect(t, checkboxEffect.MatchName, "ADBE Checkbox Control")
	expect(t, checkboxEffect.Name, "Checkbox Control")
	expect(t, checkboxEffect.PropertyType, PropertyTypeGroup)
	expect(t, len(checkboxEffect.Properties), 1)
	expect(t, checkboxEffect.Properties[0].Index, uint32(1))
	expect(t, checkboxEffect.Properties[0].MatchName, "ADBE Checkbox Control-0001")
	expect(t, checkboxEffect.Properties[0].Name, "Checkbox")
	expect(t, checkboxEffect.Properties[0].PropertyType, PropertyTypeBoolean)
	expect(t, len(checkboxEffect.Properties[0].Properties), 0)
	expect(t, len(checkboxEffect.Properties[0].SelectOptions), 0)

	sliderEffect := expressionControlsLayer.Effects[1]
	expect(t, sliderEffect.Index, uint32(2))
	expect(t, sliderEffect.MatchName, "ADBE Slider Control")
	expect(t, sliderEffect.Name, "Slider Control")
	expect(t, sliderEffect.PropertyType, PropertyTypeGroup)
	expect(t, len(sliderEffect.Properties), 1)
	expect(t, sliderEffect.Properties[0].Index, uint32(1))
	expect(t, sliderEffect.Properties[0].MatchName, "ADBE Slider Control-0001")
	expect(t, sliderEffect.Properties[0].Name, "Slider")
	expect(t, sliderEffect.Properties[0].PropertyType, PropertyTypeOneD)
	expect(t, len(sliderEffect.Properties[0].Properties), 0)
	expect(t, len(sliderEffect.Properties[0].SelectOptions), 0)

	pointEffect := expressionControlsLayer.Effects[2]
	expect(t, pointEffect.Index, uint32(3))
	expect(t, pointEffect.MatchName, "ADBE Point Control")
	expect(t, pointEffect.Name, "Point Control")
	expect(t, pointEffect.PropertyType, PropertyTypeGroup)
	expect(t, len(pointEffect.Properties), 1)
	expect(t, pointEffect.Properties[0].Index, uint32(1))
	expect(t, pointEffect.Properties[0].MatchName, "ADBE Point Control-0001")
	expect(t, pointEffect.Properties[0].Name, "Point")
	expect(t, pointEffect.Properties[0].PropertyType, PropertyTypeTwoD)
	expect(t, len(pointEffect.Properties[0].Properties), 0)
	expect(t, len(pointEffect.Properties[0].SelectOptions), 0)

	threeDPointEffect := expressionControlsLayer.Effects[3]
	expect(t, threeDPointEffect.Index, uint32(4))
	expect(t, threeDPointEffect.MatchName, "ADBE Point3D Control")
	expect(t, threeDPointEffect.Name, "3D Point Control")
	expect(t, threeDPointEffect.PropertyType, PropertyTypeGroup)
	expect(t, len(threeDPointEffect.Properties), 1)
	expect(t, threeDPointEffect.Properties[0].Index, uint32(1))
	expect(t, threeDPointEffect.Properties[0].MatchName, "ADBE Point3D Control-0001")
	expect(t, threeDPointEffect.Properties[0].Name, "3D Point")
	expect(t, threeDPointEffect.Properties[0].PropertyType, PropertyTypeThreeD)
	expect(t, len(threeDPointEffect.Properties[0].Properties), 0)
	expect(t, len(threeDPointEffect.Properties[0].SelectOptions), 0)

	colorEffect := expressionControlsLayer.Effects[4]
	expect(t, colorEffect.Index, uint32(5))
	expect(t, colorEffect.MatchName, "ADBE Color Control")
	expect(t, colorEffect.Name, "Color Control")
	expect(t, colorEffect.PropertyType, PropertyTypeGroup)
	expect(t, len(colorEffect.Properties), 1)
	expect(t, colorEffect.Properties[0].Index, uint32(1))
	expect(t, colorEffect.Properties[0].MatchName, "ADBE Color Control-0001")
	expect(t, colorEffect.Properties[0].Name, "Color")
	expect(t, colorEffect.Properties[0].PropertyType, PropertyTypeColor)
	expect(t, len(colorEffect.Properties[0].Properties), 0)
	expect(t, len(colorEffect.Properties[0].SelectOptions), 0)

	angleEffect := expressionControlsLayer.Effects[5]
	expect(t, angleEffect.Index, uint32(6))
	expect(t, angleEffect.MatchName, "ADBE Angle Control")
	expect(t, angleEffect.Name, "Angle Control")
	expect(t, angleEffect.PropertyType, PropertyTypeGroup)
	expect(t, len(angleEffect.Properties), 1)
	expect(t, angleEffect.Properties[0].Index, uint32(1))
	expect(t, angleEffect.Properties[0].MatchName, "ADBE Angle Control-0001")
	expect(t, angleEffect.Properties[0].Name, "Angle")
	expect(t, angleEffect.Properties[0].PropertyType, PropertyTypeAngle)
	expect(t, len(angleEffect.Properties[0].Properties), 0)
	expect(t, len(angleEffect.Properties[0].SelectOptions), 0)

	layerSelectEffect := expressionControlsLayer.Effects[6]
	expect(t, layerSelectEffect.Index, uint32(7))
	expect(t, layerSelectEffect.MatchName, "ADBE Layer Control")
	expect(t, layerSelectEffect.Name, "Layer Control")
	expect(t, layerSelectEffect.PropertyType, PropertyTypeGroup)
	expect(t, len(layerSelectEffect.Properties), 1)
	expect(t, layerSelectEffect.Properties[0].Index, uint32(1))
	expect(t, layerSelectEffect.Properties[0].MatchName, "ADBE Layer Control-0001")
	expect(t, layerSelectEffect.Properties[0].Name, "Layer")
	expect(t, layerSelectEffect.Properties[0].PropertyType, PropertyTypeLayerSelect)
	expect(t, len(layerSelectEffect.Properties[0].Properties), 0)
	expect(t, len(layerSelectEffect.Properties[0].SelectOptions), 0)
}
