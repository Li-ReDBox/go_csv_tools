package csv

import (
	"encoding/csv"
	"fmt"
	"log"
	"strings"
)

func Read() [][]string {
	in := `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"
`
	r := csv.NewReader(strings.NewReader(in))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(records)
	return records
}

func CreateRecords() []map[string]string {
	var records []map[string]string
	lines := Read()
	row := make(map[string]string)

	for i := 1; i < len(lines); i++ {
		for j, c := range lines[i] {
			row[lines[0][j]] = c
		}
		records = append(records, row)
	}

	return records
}

type Title map[string]int

// CreateTitle returns a Title from the input of a slice of string.
// The value of each entry is the zero-base column number from a csv file.
func CreateTitle(names []string) Title {
	t := make(Title)
	for i, n := range names {
		t[n] = i
	}
	return t
}
