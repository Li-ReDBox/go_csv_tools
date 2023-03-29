package csv

import (
	"regexp"
	"sort"
	"strconv"
)

const REG_INT = `^(\s*)(0|[1-9]\d*)(\s*)$`

var validInt = regexp.MustCompile(REG_INT)

type comparable interface {
	~string | ~int | ~float64
}

// compare returns -1 to indicate p is less then q, 1 to indicate p is greater than q,
// 0 to indicate they are equal.
func compare[C comparable](p, q C) int {
	if p < q {
		return -1
	} else if p > q {
		return 1
	}

	return 0
}

type Rows [][]string

func OrderByColumns(rows Rows, markers []int) *rowsSorter {
	return &rowsSorter{
		rows:       rows,
		markers:    markers,
		intColumns: make(map[int]struct{}),
	}
}

// rowsSorter implements the Sort interface, sorting the changes within.
type rowsSorter struct {
	rows    Rows
	markers []int
	// a tracker of int columns
	intColumns map[int]struct{}
}

// Len is part of sort.Interface.
func (byCols *rowsSorter) Len() int {
	return len(byCols.rows)
}

// Swap is part of sort.Interface.
func (byCols *rowsSorter) Swap(i, j int) {
	byCols.rows[i], byCols.rows[j] = byCols.rows[j], byCols.rows[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// marked columns until it finds a comparison that discriminates between
// the two items (one is less than the other).
// There can be multiple markers, their priorities are defined by the index in the slice.
// Only when a higher priority marker cannot make a discrimination, it passes on to the next marker.
func (byCols *rowsSorter) Less(i, j int) bool {
	// all markers need to be less
	// fmt.Println(rows[i], "vs", rows[j])
	var order, marker int

	// Check first markers, if i equals to j on the marker, continue to the next marker
	for m := 0; m < len(byCols.markers)-1; m++ {
		marker = byCols.markers[m]

		// fmt.Printf("Compare marker %d, check %s < %s\n", marker, rows[i][marker], rows[j][marker])
		if _, exists := byCols.intColumns[marker]; exists {
			// fmt.Printf("%d has been checked before\n", marker)
			p, _ := strconv.Atoi(byCols.rows[i][marker])
			q, _ := strconv.Atoi(byCols.rows[j][marker])
			order = compare(p, q)
		} else if validInt.MatchString(byCols.rows[i][marker]) && validInt.MatchString(byCols.rows[j][marker]) {
			// fmt.Printf("Adding %d as int column\n", marker)
			byCols.intColumns[marker] = struct{}{}
			p, _ := strconv.Atoi(byCols.rows[i][marker])
			q, _ := strconv.Atoi(byCols.rows[j][marker])
			order = compare(p, q)
		} else {
			order = compare(byCols.rows[i][marker], byCols.rows[j][marker])
		}

		if order != 0 {
			switch order {
			case -1:
				return true
			case 1:
				return false
			}
		}
	}

	// i and j are equal on all previous markers, with the last marker we just need to check if i is really less than j
	marker = byCols.markers[len(byCols.markers)-1]
	// fmt.Printf("Compare the last marker %d, check '%s' < '%s', result = %t\n", marker, rows[i][marker], rows[j][marker], rows[i][marker] < rows[j][marker])
	if _, exists := byCols.intColumns[marker]; exists {
		// fmt.Printf("last checker: %d has been checked before\n", marker)
		p, _ := strconv.Atoi(byCols.rows[i][marker])
		q, _ := strconv.Atoi(byCols.rows[j][marker])
		return p < q
	} else if validInt.MatchString(byCols.rows[i][marker]) && validInt.MatchString(byCols.rows[j][marker]) {
		// fmt.Printf("last checker: Adding %d as int column\n", marker)
		byCols.intColumns[marker] = struct{}{}
		p, _ := strconv.Atoi(byCols.rows[i][marker])
		q, _ := strconv.Atoi(byCols.rows[j][marker])
		return p < q
	} else {
		return byCols.rows[i][marker] < byCols.rows[j][marker]
	}
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (byCols *rowsSorter) Sort() {
	sort.Sort(byCols)
}
