package csv

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/exp/maps"
)

func basicData() string {
	return `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"
`
}

func TestRead(t *testing.T) {
	records := read(strings.NewReader(basicData()))
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

func basicRows() [][]string {
	return [][]string{
		{"first_name", "last_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}
}

func numbersAsStrings() [][]string {
	return [][]string{{"gri", "Go", "100"},
		{"ken", "C", "150"},
		{"glenda", "Go", "200"},
		{"rsc", "Go", "200"},
		{"r", "Go", "100"},
		{"ken", "Go", "200"},
		{"dmr", "C", "100"},
		{"r", "C", "150"},
		{"gri", "Smalltalk", "80"},
	}
}
func TestCreateRecords(t *testing.T) {
	rows := basicRows()

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

func TestTitle_sortingMarkers(t *testing.T) {
	input := Title{"a": 0, "b": 1, "c": 2}

	type args struct {
		nm []NamedMarker
	}
	tests := []struct {
		name string
		tr   Title
		args args
		want []Marker
	}{
		{
			"All presented", input, args{[]NamedMarker{{"a", Ascending}, {"b", Ascending}, {"c", Descending}}}, []Marker{{0, Ascending}, {1, Ascending}, {2, Descending}},
		},
		{
			"Wrong name", input, args{[]NamedMarker{{"out of range", Ascending}, {"b", Ascending}}}, []Marker{{}, {1, Ascending}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.sortingMarkers(tt.args.nm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Title.sortingMarkers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTitle_index(t *testing.T) {
	type args struct {
		names []string
	}

	input := Title{"a": 0, "b": 1, "c": 2}
	tests := []struct {
		name string
		tr   Title
		args args
		want []int
	}{
		{
			"All presented", input, args{[]string{"a", "b", "c"}}, []int{0, 1, 2},
		},
		{
			"Wrong name", input, args{[]string{"out of range"}}, []int{-1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.index(tt.args.names); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Title.index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sum(in [][]int8, c int) int8 {
	var total int8 = 0
	for _, v := range in {
		total += v[c]
	}
	fmt.Printf("sum of column %d = %d\n", c, total)
	return total
}

func TestReorder(t *testing.T) {
	names := [...]string{"first_name", "last_name", "username"}

	var count int8 = 5
	var checker int8

	data := [][]int8{
		{0, 1, 2},
		{0, 1, 2},
		{0, 1, 2},
		{0, 1, 2},
		{0, 1, 2},
	}

	titles := createTitle(names[:])

	for _, o := range [...]int{2, 0, 1} {
		if titles[names[o]] != o {
			t.Errorf("Order is %d, want %d\n", titles[names[o]], o)
		}
		// this is not the way to slice a slice of a slice, data[o][:] or data[:][o] equals to data[o]
		// so this test always fails
		checker = sum(data, o) / count
		if checker != int8(o) {
			t.Errorf("checker is %d, want %d\n", checker, o)
		}
	}
}

func ExampleProcessor_Sort_ascending() {
	markers := []Marker{{0, Ascending}, {2, Ascending}}

	p := &Processor{createTitle([]string{"user", "sub", "scores"}), numbersAsStrings()}
	p.Sort(markers)
	p.Print()

	// Output:
	// Titles:
	// map[scores:2 sub:1 user:0]
	// Rows:
	// 1, [dmr C 100]
	// 2, [glenda Go 200]
	// 3, [gri Smalltalk 80]
	// 4, [gri Go 100]
	// 5, [ken C 150]
	// 6, [ken Go 200]
	// 7, [r Go 100]
	// 8, [r C 150]
	// 9, [rsc Go 200]
}

func ExampleProcessor_Sort_mixed() {
	titles := createTitle([]string{"user", "sub", "scores"})
	nms := []NamedMarker{{"user", Ascending}, {"scores", Descending}}

	p := &Processor{titles, numbersAsStrings()}
	p.Sort(titles.sortingMarkers(nms))
	p.Print()

	// Output:
	// Titles:
	// map[scores:2 sub:1 user:0]
	// Rows:
	// 1, [dmr C 100]
	// 2, [glenda Go 200]
	// 3, [gri Go 100]
	// 4, [gri Smalltalk 80]
	// 5, [ken Go 200]
	// 6, [ken C 150]
	// 7, [r C 150]
	// 8, [r Go 100]
	// 9, [rsc Go 200]
}
