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
	if jsProj.ExpressionEngine != "javascript-1.0" {
		t.Fatalf("expected %s to equal javascript-1.0", jsProj.ExpressionEngine)
	}

	esProj, err := Open("data/ExEn-es.aep")
	if err != nil {
		t.Fatal(err)
	}
	if esProj.ExpressionEngine != "extendscript" {
		t.Fatalf("expected %s to equal extendscript", esProj.ExpressionEngine)
	}
}

// TestBitDepth tests valid parsing of a project's bit depth setting
func TestBitDepth(t *testing.T) {
	bpc32Proj, err := Open("data/BPC-32.aep")
	if err != nil {
		t.Fatal(err)
	}
	if bpc32Proj.Depth != BPC32 {
		t.Fatalf("expected %d to equal %d", bpc32Proj.Depth, BPC32)
	}
	bpc16Proj, err := Open("data/BPC-16.aep")
	if err != nil {
		t.Fatal(err)
	}
	if bpc16Proj.Depth != BPC16 {
		t.Fatalf("expected %d to equal %d", bpc16Proj.Depth, BPC16)
	}
	bpc8Proj, err := Open("data/BPC-8.aep")
	if err != nil {
		t.Fatal(err)
	}
	if bpc8Proj.Depth != BPC8 {
		t.Fatalf("expected %d to equal %d", bpc8Proj.Depth, BPC8)
	}
}
