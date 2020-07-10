package aep

import "testing"

func expect(t *testing.T, values ...interface{}) {
	if len(values) >= 1 {
		if res, ok := values[0].(bool); ok {
			if !res {
				if len(values) > 1 {
					t.Fatal(values[1])
				} else {
					t.Fatalf("Expected true value\n")
				}
			}
		} else {
			if values[0] != values[1] {
				if len(values) > 2 {
					t.Fatal(values[2])
				} else {
					t.Fatalf("Expected %+v to equal %+v\n", values[0], values[1])
				}
			}
		}
	} else {
		t.Fatalf("Invalid number of operands to expect")
	}
}
