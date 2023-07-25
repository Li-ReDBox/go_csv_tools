package csv

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// basic content for testing basic csv processes
const basicContent = `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"
`

func TestRead(t *testing.T) {
	records, _ := read(strings.NewReader(basicContent))
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
	p := &Processor{rows: numbersAsStrings()}
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

func TestProcessor_Split(t *testing.T) {
	titles := createTitle([]string{"language", "level", "value"})

	source := [][]string{
		{"C", "L1", "1"},
		{"C", "L2", "2"},
		{"C", "L1", "3"},
		{"C", "L2", "4"},
		{"C", "L1", "5"},
		{"JS", "L2", "6"},
		{"JS", "L1", "7"},
		{"JS", "L2", "8"},
		{"JS", "L1", "9"},
		{"Go", "L1", "10"},
		{"Go", "L1", "11"},
		{"Smalltalk", "L1", "12"},
	}
	p := &Processor{titles, source}

	names := []string{"level", "language"}
	inds, _ := titles.indexes(names)
	markers := []Marker{{inds[0], Descending}, {inds[1], Ascending}}
	p.Sort(markers)

	np, _ := p.Split(names)
	if len(np) != 6 {
		t.Errorf("Wanted to get 6 new Processors, but got %d", len(np))
		p.Print()
		for _, n := range np {
			n.Print()
		}
	}
}

func TestProcessor_Filter(t *testing.T) {
	var current, after runtime.MemStats

	source := [][]string{
		{"C", "L1", "1"},
		{"C", "L2", "2"}, // cl2
		{"C", "L1", "3"},
		{"C", "L2", "4"}, // cl2
		{"C", "L1", "5"},
		{"JS", "L2", "6"},
		{"JS", "L1", "7"},
		{"JS", "L2", "8"},
		{"JS", "L1", "9"},
		{"Go", "L1", "10"},
		{"Go", "L1", "11"},
		{"Smalltalk", "L1", "12"}, // smalltalk
	}

	p := &Processor{rows: source}

	cl2 := func(r []string) bool {
		if r[0] == "C" && r[1] == "L2" {
			return true
		}
		return false
	}

	smalltalk := func(r []string) bool {
		return r[0] == "Smalltalk"
	}

	runtime.ReadMemStats(&current)
	p.Filter([]Isfunc{cl2, smalltalk})

	want, nRows := 3, len(p.rows)
	if nRows != want {
		t.Errorf("Expecting rows = %d, but got %d\n", want, nRows)
	}

	runtime.GC()
	runtime.ReadMemStats(&after)
	// fmt.Printf("Before:\n %#v\n, After:\n %#v\n", current, after)
	fmt.Printf("Before:\n alloc=%v, heapalloc=%v, HeapIdle=%v, HeapReleased=%v, HeapObjects=%v, NumGC=%v\n",
		current.Alloc, current.HeapAlloc, current.HeapIdle, current.HeapReleased, current.HeapObjects, current.NumGC)
	fmt.Printf("Before:\n alloc=%v, heapalloc=%v, HeapIdle=%v, HeapReleased=%v, HeapObjects=%v, NumGC=%v\n",
		after.Alloc, after.HeapAlloc, after.HeapIdle, after.HeapReleased, after.HeapObjects, after.NumGC)
}

func Test_md5hash(t *testing.T) {
	o := []string{"These pretzels are making me thirsty."}

	if o[0] != strings.Join(o, "") {
		t.Errorf("o cannot be joined to represent o[0]: %s\n", strings.Join(o, ""))
	}

	want := "b0804ec967f48520697662a204f5fe72"

	r := md5hash(o)

	if r != want {
		t.Errorf("mdhash failed. wanted %s, but got %s\n", want, r)
	}
}

func TestProcessor_Unique(t *testing.T) {
	source := [][]string{
		{"C", "L1", "1"},
		{"JS", "L2", "6"},
		{"Go", "L1", "10"},
		{"Smalltalk", "L1", "12"},
		{"C", "L1", "1"},
		{"JS", "L2", "6"},
		{"Go", "L1", "10"},
		{"Smalltalk", "L1", "12"},
		{"C", "L1", "1"},
		{"JS", "L2", "6"},
		{"Go", "L1", "10"},
		{"Smalltalk", "L1", "12"},
	}

	p := &Processor{rows: source}
	np := p.Unique()

	if len(np.rows) != 4 {
		t.Errorf("The rows should have 4 unique rows, but np = %d\n", len(np.rows))
	}

	np.rows[0][0] = "corrupted"
	if np.rows[0][0] == p.rows[0][0] {
		t.Error("Should have copied values and created two completely separated rows, but [0][0] are equal.")
	}
}

func TestProcessorClone(t *testing.T) {
	data := basicRows()
	source := &Processor{
		createTitle(data[0]),
		data[1:],
	}

	copy := source.Clone()

	if copy == source {
		t.Errorf("Clone did not produce a new Processor, want = %v, got = %v\n", source, copy)
	}

	if !maps.Equal(source.titles, copy.titles) {
		t.Errorf("Clone failed to reproduce titles: want = %v, got = %v\n", source.titles, copy.titles)
	}

	// This is a silly test
	copy.titles["extra"] = len(source.titles) + 1

	if maps.Equal(source.titles, copy.titles) {
		t.Errorf("Clone failed to create a new copy: want = %v != %v\n", source.titles, copy.titles)
	}

	nRows := len(source.rows)
	for i := 0; i < nRows; i++ {
		copy.rows[i][0] = ""
		if copy.rows[i][0] == source.rows[i][0] {
			t.Error("Rows are not fully cloned: changed copy went to source.")
		}
	}
}

func TestReplaceRows(t *testing.T) {
	p := &Processor{
		rows: numbersAsStrings()}

	suspectsWithHighScore := func(elems []string) bool {
		// those suspects cannot achieve higher than 100 in any course. If this is the case, mark them
		suspects := []string{"glenda", "ken", "baz"}
		if slices.Contains(suspects, elems[0]) {
			score, err := strconv.Atoi(elems[2])
			if err == nil && score > 100 {
				return true
			}
		}
		return false
	}

	mark := func(s []string) {
		s[2] = "Got you"
	}

	op := Operation{
		Check: suspectsWithHighScore,
		Act:   mark,
	}
	p.Replace([]Operation{op})

	marked := func(elems []string) bool {
		return elems[2] == "Got you"
	}

	// check how many suspects are caught
	p.Filter([]Isfunc{marked})
	if len(p.rows) == 0 {
		t.Error("Replace failed because no rows contains 'Got you' is found")
		p.Print()
	}
}

func TestDerive(t *testing.T) {
	records, _ := read(strings.NewReader(basicContent))
	p := &Processor{createTitle(records[0][:]), records[1:][:]}

	ban := func(fName, lName string) string {
		return fmt.Sprintf("%s %s has been banned", fName, lName)
	}

	p.Derive(0, 1, "Message", ban)

	for _, row := range p.rows {
		if row[3] != ban(row[0], row[1]) {
			t.Errorf("Failed to ban %s %s. Message is: %s", row[0], row[1], row[3])
		}
	}
}
