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

func TestReplace(t *testing.T) {
	original := []string{"a", "b", "c"}
	markers := []int{0, 2}

	prefix := func(content []string) {
		for i := 0; i < len(content); i++ {
			content[i] = "prefix-" + content[i]
		}
	}

	replace(original, markers, prefix)

	for _, i := range markers {
		if !strings.HasPrefix(original[i], "prefix-") {
			t.Error("Expecting to have prefix: prefix-, but had: ", original[i])
		}
	}
}
