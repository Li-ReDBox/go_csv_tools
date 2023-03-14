package csv

import (
	"fmt"
	"testing"
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
