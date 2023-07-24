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
func ExampleProcessor_Split() {
	titles := createTitle([]string{"user", "sub", "scores"})
	p := &Processor{titles, numbersAsStrings()}

	names := []string{"sub"}
	inds, _ := titles.indexes(names)
	markers := []Marker{{inds[0], Descending}}
	p.Sort(markers)
	fmt.Println("Source:")
	p.Print()

	fmt.Println("Split results:")
	np, _ := p.Split(names)
	for _, s := range np {
		s.Print()
	}

	// Output:
	// Source:
	// Titles:
	// user, sub, scores
	// Rows:
	// 1 gri, Smalltalk, 80
	// 2 gri, Go, 100
	// 3 glenda, Go, 200
	// 4 rsc, Go, 200
	// 5 r, Go, 100
	// 6 ken, Go, 200
	// 7 ken, C, 150
	// 8 dmr, C, 100
	// 9 r, C, 150
	//
	// Split results:
	// Titles:
	// user, sub, scores
	// Rows:
	// 1 gri, Smalltalk, 80
	//
	// Titles:
	// user, sub, scores
	// Rows:
	// 1 gri, Go, 100
	// 2 glenda, Go, 200
	// 3 rsc, Go, 200
	// 4 r, Go, 100
	// 5 ken, Go, 200
	//
	// Titles:
	// user, sub, scores
	// Rows:
	// 1 ken, C, 150
	// 2 dmr, C, 100
	// 3 r, C, 150
	//
}

func ExampleProcessor_Unique() {
	dup := append(basicRows(), basicRows()...)
	p := Processor{rows: dup}
	np := p.Unique()
	np.Print()

	// Output:
	// Titles:
	//
	// Rows:
	// 1 first_name, last_name, username
	// 2 Rob, Pike, rob
	// 3 Ken, Thompson, ken
	// 4 Robert, Griesemer, gri
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

// Example_sortDates demonstrates how to formalise dates to ISO 8601 and use it in sorting.
// Note: If replacing to ISO 8601 dates is not desirable, sorting can be done by duplicating a column of date (not implemented)
// for sorting then remove var Convert method.
// Or var two replacements: first to ISO 8601 then replace the column back after sorting. No new feature is needed.
func Example_sortDates() {
	const dates = `date
12/12/2005
31/1/2005
1/1/2005
`
	records := read(strings.NewReader(dates))
	p := &Processor{createTitle(records[0][:]), records[1:][:]}

	pad := func(d string) string {
		// only day and month are processed in this example
		if len(d) == 2 {
			return d
		}
		return "0" + d
	}

	formatter := func(elems []string) {
		ps := strings.Split(elems[0], "/")

		// reformat to ISO 8601
		elems[0] = fmt.Sprintf("%s-%s-%s", ps[2], pad(ps[1]), pad(ps[0]))
	}

	op := Operation{
		Check: func(elems []string) bool { return true },
		Act:   formatter,
	}
	p.Replace([]Operation{op})

	markers := []Marker{{0, Descending}}
	p.Sort(markers)

	p.Print()

	// Output:
	// Titles:
	// date
	// Rows:
	// 1 2005-12-12
	// 2 2005-01-31
	// 3 2005-01-01
}
