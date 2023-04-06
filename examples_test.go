package csv

import (
	"fmt"
	"strings"
)

func ExampleProcessor_Swap() {
	p := &Processor{createTitle([]string{"user", "sub", "scores"}), numbersAsStrings()}
	p.Swap("sub", "user")
	p.Print()

	// Output:
	// Titles:
	// sub, user, scores
	// Rows:
	// 1 Go, gri, 100
	// 2 C, ken, 150
	// 3 Go, glenda, 200
	// 4 Go, rsc, 200
	// 5 Go, r, 100
	// 6 Go, ken, 200
	// 7 C, dmr, 100
	// 8 C, r, 150
	// 9 Smalltalk, gri, 80
}

func ExampleProcessor_Sort_ascending() {
	markers := []Marker{{0, Ascending}, {2, Ascending}}

	p := &Processor{createTitle([]string{"user", "sub", "scores"}), numbersAsStrings()}
	p.Sort(markers)
	p.Print()

	// Output:
	// Titles:
	// user, sub, scores
	// Rows:
	// 1 dmr, C, 100
	// 2 glenda, Go, 200
	// 3 gri, Smalltalk, 80
	// 4 gri, Go, 100
	// 5 ken, C, 150
	// 6 ken, Go, 200
	// 7 r, Go, 100
	// 8 r, C, 150
	// 9 rsc, Go, 200
}

func ExampleProcessor_Sort_mixed() {
	titles := createTitle([]string{"user", "sub", "scores"})
	nms := []NamedMarker{{"user", Ascending}, {"scores", Descending}}

	p := &Processor{titles, numbersAsStrings()}
	markers, err := titles.sortingMarkers(nms)
	if err == nil {
		p.Sort(markers)
		p.Print()
	}

	// Output:
	// Titles:
	// user, sub, scores
	// Rows:
	// 1 dmr, C, 100
	// 2 glenda, Go, 200
	// 3 gri, Go, 100
	// 4 gri, Smalltalk, 80
	// 5 ken, Go, 200
	// 6 ken, C, 150
	// 7 r, C, 150
	// 8 r, Go, 100
	// 9 rsc, Go, 200
}

func ExampleProcessor_Extract() {
	p := &Processor{createTitle([]string{"user", "sub", "scores"}), numbersAsStrings()}
	sub, _ := p.Extract([]string{"sub", "user"})

	for i, row := range sub {
		fmt.Println(i+1, strings.Join(row, ", "))
	}

	// Output:
	// 1 Go, gri
	// 2 C, ken
	// 3 Go, glenda
	// 4 Go, rsc
	// 5 Go, r
	// 6 Go, ken
	// 7 C, dmr
	// 8 C, r
	// 9 Smalltalk, gri
}

func ExampleProcessor_Convert() {
	p := &Processor{createTitle([]string{"user", "sub", "scores"}), numbersAsStrings()}
	c, err := p.Convert([]string{"scores", "user"})

	if err == nil {
		// mapping new titles into a slice of sorting markers
		markers := []Marker{{0, Descending}, {1, Ascending}}
		c.Sort(markers)
		c.Print()
	}

	// Output:
	// Titles:
	// scores, user
	// Rows:
	// 1 200, glenda
	// 2 200, ken
	// 3 200, rsc
	// 4 150, ken
	// 5 150, r
	// 6 100, dmr
	// 7 100, gri
	// 8 100, r
	// 9 80, gri
}

// ExampleProcessor_split shows how to split a sorted dataset
func ExampleProcessor_split() {
	titles := createTitle([]string{"user", "sub", "scores"})
	p := &Processor{titles, numbersAsStrings()}

	inds, _ := titles.indexes([]string{"user", "scores"})
	markers := []Marker{{inds[1], Ascending}, {inds[0], Ascending}}
	p.Sort(markers)

	// say we want to split by "scores": inds[1]

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

// Examples for rows
func Example_sortRows() {
	rows := numbersAsStrings()
	markers := []Marker{{0, Ascending}, {2, Ascending}}
	sorter := OrderByColumns(markers)
	sorter.Sort(rows)

	fmt.Println(rows)
	// Output:
	// [[dmr C 100] [glenda Go 200] [gri Smalltalk 80] [gri Go 100] [ken C 150] [ken Go 200] [r Go 100] [r C 150] [rsc Go 200]]
}
