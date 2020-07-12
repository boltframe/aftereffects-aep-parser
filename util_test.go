package aep

import (
	"runtime"
	"strings"
	"testing"
)

func expect(t *testing.T, values ...interface{}) {
	_, path, line, _ := runtime.Caller(1)
	file := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
	if len(values) >= 1 {
		if res, ok := values[0].(bool); ok {
			if !res {
				if len(values) > 1 {
					t.Fatal(values[1])
				} else {
					t.Fatalf("(%s:%d) Expected true value\n", file, line)
				}
			}
		} else {
			if values[0] != values[1] {
				if len(values) > 2 {
					t.Fatal(values[2])
				} else {
					t.Fatalf("(%s:%d) Expected \"%+v\" to equal \"%+v\"\n", file, line, values[0], values[1])
				}
			}
		}
	} else {
		t.Fatalf("Invalid number of operands to expect")
	}
}
