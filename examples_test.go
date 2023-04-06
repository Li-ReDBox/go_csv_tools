package csv

import "fmt"

// ExampleProcessor_split shows how to split a sorted dataset
func ExampleProcessor_split() {
	titles := createTitle([]string{"user", "sub", "scores"})
	p := &Processor{titles, numbersAsStrings()}

	inds, _ := titles.indexes([]string{"user", "scores"})
	markers := []Marker{{inds[1], Ascending}, {inds[0], Ascending}}
	p.Sort(markers)

	// say we want to split by "user": inds[0]

	current := ""
	for _, r := range p.rows {
		if r[inds[1]] != current {
			current = r[inds[1]]
			fmt.Println("Section", current)
		}
		fmt.Println(r)
	}

	// Output:
	// Section 80
	// [gri Smalltalk 80]
	// Section 100
	// [dmr C 100]
	// [gri Go 100]
	// [r Go 100]
	// Section 150
	// [ken C 150]
	// [r C 150]
	// Section 200
	// [glenda Go 200]
	// [ken Go 200]
	// [rsc Go 200]
}
