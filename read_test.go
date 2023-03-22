package csv

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/exp/maps"
)

func TestRead(t *testing.T) {
	in := `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"
`

	records := read(strings.NewReader(in))
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
	rows := [][]string{
		{"first_name", "last_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}

	want := [...]map[string]string{
		{"first_name": "Rob",
			"last_name": "Pike",
			"username":  "rob",
		},
		{"first_name": "Ken",
			"last_name": "Thompson",
			"username":  "ken",
		},
		{"first_name": "Robert",
			"last_name": "Griesemer",
			"username":  "gri",
		},
	}
	records := createRecords(rows)
	if len(records) != 3 {
		t.Errorf("Want size to be greater 3, but %d", len(records))
	}

	for i, r := range records {
		if !maps.Equal(r, want[i]) {
			t.Errorf("createRecords() element %d is %v, want %v", i, r, want[i])
		}
	}
}

func TestCreateTitle(t *testing.T) {
	names := [...]string{"a", "b", "c"}

	want := Title{"a": 0, "b": 1, "c": 2}

	titles := createTitle(names[:])

	if !maps.Equal(titles, want) {
		t.Errorf("CreateTitle() = %v, want %v", titles, want)
	}
}
