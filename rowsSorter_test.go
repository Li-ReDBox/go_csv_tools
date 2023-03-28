package csv

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"testing"
)

const REG_INT = `^(\s*)(0|[1-9]\d*)(\s*)$`

var validInt = regexp.MustCompile(REG_INT)

type Rows [][]string

// Len is part of sort.Interface.
func (rows Rows) Len() int {
	return len(rows)
}

// Swap is part of sort.Interface.
func (rows Rows) Swap(i, j int) {
	rows[i], rows[j] = rows[j], rows[i]
}

// slices of int type data columns, the indicators is not enough
// the key is the column index of original data, the value is the index of intColumns
var intTypes = make(map[int]int)

type intColumns [][]int

func (cols intColumns) less(n, i, j int) int {
	p, q := cols[n][i], cols[n][j]
	switch {
	case p < q:
		// p < q, so we have a decision.
		return -1
	case p > q:
		// p > q, so we have a decision.
		return 1
	}
	return 0
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.

// This is a multiple key comparison, its priorities are defined by the order to markers.
// Only the higher priority marker cannot make an decision, it passes on to the next marker.
func (rows Rows) Less(i, j int) bool {
	// which columns are used in comparison
	// all markers need to be less
	// fmt.Println(rows[i], "vs", rows[j])
	markers := []int{0, 2}
	var m, order int
	for m = 0; m < len(markers)-1; m++ {
		// fmt.Printf("Compare marker %d, check %s < %s\n", markers[m], rows[i][markers[m]], rows[j][markers[m]])
		if validInt.MatchString(rows[i][markers[m]]) && validInt.MatchString(rows[j][markers[m]]) {
			p, _ := strconv.Atoi(rows[i][markers[m]])
			q, _ := strconv.Atoi(rows[j][markers[m]])
			order = compare(p, q)
		} else {
			order = compare(rows[i][markers[m]], rows[j][markers[m]])
		}

		if order != 0 {
			switch order {
			case -1:
				return true
			case 1:
				return false
			}
		}
		// switch {
		// case rows[i][markers[m]] < rows[j][markers[m]]:
		// 	// p < q, so we have a decision.
		// 	return true
		// case rows[i][markers[m]] > rows[j][markers[m]]:
		// 	// p > q, so we have a decision.
		// 	return false
		// }
	}
	// fmt.Printf("Compare the last marker %d, check '%s' < '%s', result = %t\n", markers[m], rows[i][markers[m]], rows[j][markers[m]], rows[i][markers[m]] < rows[j][markers[m]])
	if validInt.MatchString(rows[i][markers[m]]) && validInt.MatchString(rows[j][markers[m]]) {
		p, _ := strconv.Atoi(rows[i][markers[m]])
		q, _ := strconv.Atoi(rows[j][markers[m]])
		return p < q
	} else {
		return rows[i][markers[m]] < rows[j][markers[m]]
	}
}

func Example_sortRows() {
	var rows Rows = Rows{
		{"gri", "Go", "100"},
		{"ken", "C", "150"},
		{"glenda", "Go", "200"},
		{"rsc", "Go", "200"},
		{"r", "Go", "100"},
		{"ken", "Go", "200"},
		{"dmr", "C", "100"},
		{"r", "C", "150"},
		{"gri", "Smalltalk", "80"},
	}

	sort.Sort(&rows)

	fmt.Println(rows)
	// Output:
	// [[dmr C 100] [glenda Go 200] [gri Smalltalk 80] [gri Go 100] [ken C 150] [ken Go 200] [r Go 100] [r C 150] [rsc Go 200]]
}

/*
excel sorting:
dmr	C	100
glenda	Go	200
gri	Smalltalk	80
gri	Go	100
ken	C	150
ken	Go	200
r	Go	100
r	C	150
rsc	Go	200

*/

func TestCompareNumStrings(t *testing.T) {
	tests := []struct {
		s1, s2 string
		r      bool
	}{
		{"", "", false},
		{"80", "100", false}, // this string comparison is different to numbers
		{"80", "81", true},
		{"7", "81000", true},  // this string comparison is different to numbers
		{"8", "81000", true},  // this string comparison is different to numbers
		{"9", "81000", false}, // this string comparison is different to numbers
		{"a", "b", true},
		{"b", "a", false},
		{"a", "aa", true},
		{"aa", "a", false},
	}
	for _, s := range tests {
		cmp := s.s1 < s.s2
		if cmp != s.r {
			t.Errorf("%s < %s failed, it is %t, but want %t", s.s1, s.s2, cmp, s.r)
		}
	}
}

func TestIntType(t *testing.T) {
	// Compile the expression once, usually at init time.
	// Use raw strings to avoid having to quote the backslashes.

	numbers := []string{"1", "0", " 1", " 0", "1  ", "0  ", " 1 ", " 0 ", "150", "80", "201", " 487576  "}
	for _, n := range numbers {
		if !validInt.MatchString(n) {
			t.Errorf("%s should match to an int", n)
		}
		if validInt.MatchString(n) {
			fmt.Printf("Extracted int from %q = %q \n", n, validInt.FindStringSubmatch(n)[2])
		}
	}

	nonNumbers := []string{"010", "09", " 010", " 010", "010  "}
	for _, n := range nonNumbers {
		if validInt.MatchString(n) {
			t.Errorf("%s should not match to an int", n)
		}
	}
}

func TestFindInts(t *testing.T) {
	// mixes is a matrix, so each row has the same cells
	mixes := numbersAsStrings()

	nLines := len(mixes)
	nCells := len(mixes[0])
	var isInt = make([]bool, nCells)
	for c := 0; c < nCells; c++ {
		isInt[c] = true
		for l := 0; l < nLines; l++ {
			isInt[c] = isInt[c] && validInt.MatchString(mixes[l][c])
			// if !validInt.MatchString(mixes[l][c]) {
			// 	t.Errorf("Pattern does not think %s is int\n", mixes[l][c])
			// }
		}
	}
	if isInt[0] || isInt[1] || !isInt[2] {
		t.Errorf("Expected false, false, true, but had %v\n", isInt)
	}
}

type comparable interface {
	~string | ~int | ~float32 | ~float64
}

func compare[C comparable](p, q C) int {
	if p < q {
		return -1
	} else if p > q {
		return 1
	}

	return 0
}

func TestCompare(t *testing.T) {
	if compare(1, 2) != -1 {
		t.Errorf("int compare failed")
	}
	if compare(1.0, 2.0) != -1 {
		t.Errorf("int compare failed")
	}
	if compare("1", "2") != -1 {
		t.Errorf("int compare failed")
	}
	if compare(float64(1), float64(2)) != -1 {
		t.Errorf("int compare failed")
	}
}
