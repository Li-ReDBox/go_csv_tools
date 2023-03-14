package csv

import (
	"fmt"
	"testing"

	"golang.org/x/exp/maps"
)

func TestRead(t *testing.T) {
	records := Read()
	fmt.Println(records)

	if len(records) == 0 {
		t.Errorf("Want size to be greater than 0, but %d", len(records))
	}

	for i, v := range records[:] {
		fmt.Printf("index = %d has %v\n", i, v)
		for j, c := range records[i][:] {
			fmt.Printf("Cell %d has %s\n", j, c)
		}
		fmt.Println()
	}
}

func TestCreateRecords(t *testing.T) {
	records := CreateRecords()
	if len(records) == 0 {
		t.Errorf("Want size to be greater than 0, but %d", len(records))
	}
}

func TestCreateTitle(t *testing.T) {
	names := [...]string{"a", "b", "c"}

	want := Title{"a": 0, "b": 1, "c": 2}

	titles := CreateTitle(names[:])

	if !maps.Equal(titles, want) {
		t.Errorf("CreateTitle() = %v, want %v", titles, want)
	}
}
