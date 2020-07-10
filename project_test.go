package aep

import (
	"testing"
)

// TestExpressionEngine tests valid parsing of a project's expression engine setting
func TestExpressionEngine(t *testing.T) {
	jsProj, err := Open("data/ExEn-js.aep")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, jsProj.ExpressionEngine, "javascript-1.0")

	esProj, err := Open("data/ExEn-es.aep")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, esProj.ExpressionEngine, "extendscript")
}

// TestBitDepth tests valid parsing of a project's bit depth setting
func TestBitDepth(t *testing.T) {
	bpc32Proj, err := Open("data/BPC-32.aep")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, bpc32Proj.Depth, BPC32)

	bpc16Proj, err := Open("data/BPC-16.aep")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, bpc16Proj.Depth, BPC16)

	bpc8Proj, err := Open("data/BPC-8.aep")
	if err != nil {
		t.Fatal(err)
	}
	expect(t, bpc8Proj.Depth, BPC8)
}
