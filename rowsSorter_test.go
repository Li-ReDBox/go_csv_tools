package csv

import (
	"fmt"
	"testing"
)

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

	sorter := OrderByColumns(rows, []int{0, 2})
	sorter.Sort()

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
