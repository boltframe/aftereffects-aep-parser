package aep

import (
	"github.com/rioam2/rifx"
)

// LayerQualityLevel denotes the quality level of a layer (eg: Best, Draft, Wireframe)
type LayerQualityLevel uint16

const (
	// LayerQualityBest enumerates the value of a layer with Best Quality
	LayerQualityBest LayerQualityLevel = 0x0002
	// LayerQualityDraft enumerates the value of a layer with Draft Quality
	LayerQualityDraft LayerQualityLevel = 0x0001
	// LayerQualityWireframe enumerates the value of a layer with Wireframe Quality
	LayerQualityWireframe LayerQualityLevel = 0x0000
)

// LayerSamplingMode denotes the sampling mode of a layer (eg: Bilinear, Bicubic)
type LayerSamplingMode byte

const (
	// LayerSamplingModeBilinear enumerates the value of a layer with Bilinear Sampling
	LayerSamplingModeBilinear LayerSamplingMode = 0x00
	// LayerSamplingModeBicubic enumerates the value of a layer with Bicubic Sampling
	LayerSamplingModeBicubic LayerSamplingMode = 0x01
)

// LayerFrameBlendMode denotes the frame blending mode of a layer (eg: Frame mix, Pixel motion)
type LayerFrameBlendMode byte

const (
	// LayerFrameBlendModeFrameMix enumerates the value of a layer with Frame Mix Frame Blending
	LayerFrameBlendModeFrameMix LayerFrameBlendMode = 0x00
	// LayerFrameBlendModePixelMotion enumerates the value of a layer with Pixel Motion Frame Blending
	LayerFrameBlendModePixelMotion LayerFrameBlendMode = 0x01
)

// Layer describes a single layer in a composition.
type Layer struct {
	Index                    uint32
	Name                     string
	SourceID                 uint32
	Quality                  LayerQualityLevel
	SamplingMode             LayerSamplingMode
	FrameBlendMode           LayerFrameBlendMode
	GuideEnabled             bool
	SoloEnabled              bool
	ThreeDEnabled            bool
	AdjustmentLayerEnabled   bool
	CollapseTransformEnabled bool
	ShyEnabled               bool
	LockEnabled              bool
	FrameBlendEnabled        bool
	MotionBlurEnabled        bool
	EffectsEnabled           bool
	AudioEnabled             bool
	VideoEnabled             bool
}

func parseLayer(layerHead *rifx.List, project *Project) (*Layer, error) {
	layer := &Layer{}

	type LDTA struct {
		Unknown00     [4]byte           // Offset 0B
		Quality       LayerQualityLevel // Offset 4B
		Unknown01     [31]byte          // Offset 6B
		LayerAttrBits [3]byte           // Offset 37B
		SourceID      uint32            // Offset 40B
	}
	ldtaBlock, err := layerHead.FindByType("ldta")
	if err != nil {
		return nil, err
	}
	ldta := &LDTA{}
	ldtaBlock.ToStruct(ldta)
	layer.SourceID = ldta.SourceID
	layer.Quality = ldta.Quality
	layer.SamplingMode = LayerSamplingMode((ldta.LayerAttrBits[0] & (1 << 6)) >> 6)
	layer.FrameBlendMode = LayerFrameBlendMode((ldta.LayerAttrBits[0] & (1 << 2)) >> 2)
	layer.GuideEnabled = ((ldta.LayerAttrBits[0] & (1 << 1)) >> 1) == 1
	layer.SoloEnabled = ((ldta.LayerAttrBits[1] & (1 << 3)) >> 3) == 1
	layer.ThreeDEnabled = ((ldta.LayerAttrBits[1] & (1 << 2)) >> 2) == 1
	layer.AdjustmentLayerEnabled = ((ldta.LayerAttrBits[1] & (1 << 1)) >> 1) == 1
	layer.CollapseTransformEnabled = ((ldta.LayerAttrBits[2] & (1 << 7)) >> 7) == 1
	layer.ShyEnabled = ((ldta.LayerAttrBits[2] & (1 << 6)) >> 6) == 1
	layer.LockEnabled = ((ldta.LayerAttrBits[2] & (1 << 5)) >> 5) == 1
	layer.FrameBlendEnabled = ((ldta.LayerAttrBits[2] & (1 << 4)) >> 4) == 1
	layer.MotionBlurEnabled = ((ldta.LayerAttrBits[2] & (1 << 3)) >> 3) == 1
	layer.EffectsEnabled = ((ldta.LayerAttrBits[2] & (1 << 2)) >> 2) == 1
	layer.AudioEnabled = ((ldta.LayerAttrBits[2] & (1 << 1)) >> 1) == 1
	layer.VideoEnabled = ((ldta.LayerAttrBits[2] & (1 << 0)) >> 0) == 1

	nameBlock, err := layerHead.FindByType("Utf8")
	if err != nil {
		return nil, err
	}
	layer.Name = nameBlock.ToString()
	return layer, nil
}
