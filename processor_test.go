package csv

import (
	"errors"
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

	t.Run("All presented", func(t *testing.T) {
		want := []Marker{{0, Ascending}, {1, Ascending}, {2, Descending}}
		got, err := input.sortingMarkers([]NamedMarker{{"a", Ascending}, {"b", Ascending}, {"c", Descending}})
		if !(err == nil && reflect.DeepEqual(got, want)) {
			t.Errorf("Title.sortingMarkers() = %v, want %v", got, want)
		}
	})
	t.Run("Wrong name", func(t *testing.T) {
		got, err := input.sortingMarkers([]NamedMarker{{"out of range", Ascending}, {"b", Ascending}})
		if err != nil && got != nil {
			t.Errorf("Title.sortingMarkers() should has a non-nil error, but it is %s, markers should be nil, but %v", err, got)
		}
	})
}

func TestTitle_index(t *testing.T) {
	input := Title{"a": 0, "b": 1, "c": 2}

	t.Run("All presented", func(t *testing.T) {
		want := []int{0, 1, 2}
		got, err := input.indexes([]string{"a", "b", "c"})
		if !(err == nil && reflect.DeepEqual(got, want)) {
			t.Errorf("Title.index() = %v, want %v\n", got, want)
		}
	})
	t.Run("Wrong name", func(t *testing.T) {
		got, err := input.indexes([]string{"out of range"})
		if err != nil && got != nil {
			t.Errorf("Title.index() should has a non-nil error, but it is %s, markers should be nil, but %v", err, got)
		}
	})
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

func TestProcessor_Swap_notFound(t *testing.T) {
	p := &Processor{createTitle([]string{"user", "sub", "scores"}), numbersAsStrings()}
	err := p.Swap("first_name", "last_name")

	if err == nil {
		t.Errorf("Non-existing title should return error but received nil\n")
	}

	want := TitleNotFound("first_name")
	if err != want {
		t.Errorf("Want error as %s, but received %s\n", want, err)
	}
}

func TestTitleNotFound_equal(t *testing.T) {
	err1 := TitleNotFound("n")
	err2 := TitleNotFound("n")

	if err1 != err2 {
		t.Errorf("Same TitleNotFound error are different: %v != %v", err1, err2)
	}
}
func TestTitleNotFound_Is(t *testing.T) {
	title := "n"
	want := "csv/TitleNotFound: " + title
	var err error = TitleNotFound(title)

	if err.Error() != want {
		t.Errorf("TitleNotFound message is not well formatted")
	}

	if !errors.Is(err, TitleNotFound(title)) {
		t.Errorf("TitleNotFound cannot use Is to compare.")
	}

	err = errors.New(want)
	fmt.Println("Checking ", err)
	if errors.Is(err, TitleNotFound(title)) {
		t.Errorf("errors.Is method does not work, misclassified %s as TitleNotFound", err)
	}
}

func TestProcessor_Write(t *testing.T) {
	p := &Processor{
		rows: numbersAsStrings()}
	o := `gri,Go,100
ken,C,150
glenda,Go,200
rsc,Go,200
r,Go,100
ken,Go,200
dmr,C,100
r,C,150
gri,Smalltalk,80
`
	var w strings.Builder
	err := p.Write(&w)
	if err != nil {
		t.Errorf("Wanted err to be nil, but it is %s\n", err)
	}

	if w.String() != o {
		t.Errorf("Write method generated un expected content. wanted:\n%s\n but got:\n%s\n", o, w.String())
	}
}