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

type Direction int

const (
	Ascending  = Direction(1)
	Descending = Direction(-1)
)

// Marker defines a numbered sorting order. It can be applied to a data set of [][]string for sorting.
type Marker struct {
	Index int
	Order Direction
}

// NamedMarker defines named sorting order. It cannot be applied to a data set of [][]string for sorting.
// It needs to be mapped Marker first. Title.sortingMarkers provides a such mapping. Because it needs a mapping
// step, they are/can be validated against dataset, which is safer compared to using Marker.
type NamedMarker struct {
	Name  string
	Order Direction
}

// OrderByColumns creates a rowsSorter with a slice of Marker.
// It is anticipated to have all the markers presents in the data
// to be sorted. If any index is out of range, it will panic:
// no error handling has be defined yet.
func OrderByColumns(markers []Marker) *rowsSorter {
	return &rowsSorter{
		markers: markers,
	}
}

// rowsSorter implements the Sort interface, sorting the string rows by markers (prioritised columns)
// Currently it will check if any marked columns is int when sorting. If the column is int type, internally
// it will convert string to int then compare them in the native way.
// It support sorting in Direction: either Ascending or Descending.
type rowsSorter struct {
	rows    [][]string
	markers []Marker
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
	nMarker := len(byCols.markers)

	// Check first markers, if i equals to j on the marker, continue to the next marker
	for m := 0; m < nMarker; m++ {
		marker = byCols.markers[m].Index

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

		// apply ordering
		order = order * int(byCols.markers[m].Order)

		if order != 0 && m < nMarker-1 {
			switch order {
			case -1:
				return true
			case 1:
				return false
			}
		}
	}
	return order == -1
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (byCols *rowsSorter) Sort(rows [][]string) {
	byCols.rows = rows
	byCols.intColumns = make(map[int]struct{})

	sort.Sort(byCols)
}
