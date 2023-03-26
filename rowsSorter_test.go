package csv

import (
	"fmt"
	"regexp"
	"sort"
	"testing"
)

type Rows [][]string

// Len is part of sort.Interface.
func (rows Rows) Len() int {
	return len(rows)
}

// Swap is part of sort.Interface.
func (rows Rows) Swap(i, j int) {
	rows[i], rows[j] = rows[j], rows[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (rows Rows) Less(i, j int) bool {
	// which columns are used in comparison
	// all markers need to be less
	fmt.Println(rows[i], "vs", rows[j])
	markers := []int{0, 2}
	var m int
	for m = 0; m < len(markers)-1; m++ {
		fmt.Printf("Compare marker %d, check %s < %s\n", markers[m], rows[i][markers[m]], rows[j][markers[m]])
		switch {
		case rows[i][markers[m]] < rows[j][markers[m]]:
			// p < q, so we have a decision.
			return true
		case rows[i][markers[m]] > rows[j][markers[m]]:
			// p > q, so we have a decision.
			return false
		}
	}
	fmt.Printf("Compare the last marker %d, check '%s' < '%s', result = %t\n", markers[m], rows[i][markers[m]], rows[j][markers[m]], rows[i][markers[m]] < rows[j][markers[m]])
	return rows[i][markers[m]] < rows[j][markers[m]]
}

func TestRows(t *testing.T) {
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
	// var rows Rows = Rows{
	// 	{"z", "1", "d"},
	// 	{"b", "2", "y"},
	// 	{"c", "3", "b"},
	// 	{"c", "3", "d"},
	// }
	sort.Sort(&rows)

	if "Robert" > "first_name" {
		t.Errorf("Robert is actually > first_name which is bizarre.\n")
	}

	if !("Ken" < "Rob" && "Rob" < "Robert" && "Robert" < "first_name") {
		t.Errorf("Column 0 comparison failed, wanted %s", "Ken < Rob < Robert < first_name")
	}

	if "ken" < "rob" && "rob" < "gri" && "gri" < "username" {
		t.Errorf("Column 2 comparison failed, wanted %s", "ken < rob < gri < username")
	}

	// this is to remind me sorting of strings is different to numbers.
	if "80" > "100" {
		// '313030=100' < '3830=80' = true
		/*
					int32
			runtime·strcmp(byte *s1, byte *s2)
			{
			    uintptr i;
			    byte c1, c2;

			    for(i=0;; i++) {
			        c1 = s1[i];
			        c2 = s2[i];
			        if(c1 < c2)
			            return -1;
			        if(c1 > c2)
			            return +1;
			        if(c1 == 0)
			            return 0;
			    }
			}
		*/
		t.Errorf("I think string 80 should be less then 100, actually less check returns %t\n", "80" < "100")
	}

	fmt.Println(rows)
}

/*
sorted by 0 only
[[dmr C 100] [glenda Go 200] [gri Smalltalk 80] [gri Go 100] [ken Go 200] [ken C 150] [r C 150] [r Go 100] [rsc Go 200]]
sorted by 0 then 2
[[dmr C 100] [glenda Go 200] [gri Go 100] [gri Smalltalk 80] [ken C 150] [ken Go 200] [r C 150] [r Go 100] [rsc Go 200]]

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



Excel
b	2	y
c	3	b
c	3	d
z	1	d

[[b 2 y] [c 3 b] [c 3 d] [z 1 d]]
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
	var validInt = regexp.MustCompile(`^(\s*)(0|[1-9]+)(\s*)$`)

	numbers := []string{"1", "0", " 1", " 0", "1  ", "0  ", " 1 ", " 0 "}
	for _, n := range numbers {
		if !validInt.MatchString(n) {
			t.Errorf("%s should match to an int", n)
		}
		fmt.Printf("Extracted int from %q = %q \n", n, validInt.FindStringSubmatch(n)[2])
	}

	nonNumbers := []string{"010", "09", " 010", " 010", "010  "}
	for _, n := range nonNumbers {
		if validInt.MatchString(n) {
			t.Errorf("%s should not match to an int", n)
		}
	}
}
