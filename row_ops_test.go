package csv

import (
	"fmt"
	"strings"
	"testing"
)

func TestRowOneColOp(t *testing.T) {
	records := read(strings.NewReader(basicContent))
	p := &Processor{createTitle(records[0][:]), records[1:][:]}

	fmt.Println(p.Size())
}

func TestBiOp(t *testing.T) {
	original := []string{"a", "b", "c"}

	concat := func(a, b string) string {
		return a + b
	}
	result := biop(original, 0, 2, concat)

	if len(original)+1 != len(result) {
		t.Error("Expecting only increasing by one, got ", len(result)-len(original))
	}

	if result[len(result)-1] != original[0]+original[2] {
		t.Errorf("Expecting the last element to be %s, but got %s\n", original[0]+original[2], result[len(result)-1])
	}
}
