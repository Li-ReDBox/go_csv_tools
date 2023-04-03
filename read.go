package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const titleNotFoundPrefix = "csv/TitleNotFound"

// TitleMisMatchError describes an error when a user provided title cannot be found in a given Title.
type TitleNotFound string

func (e TitleNotFound) Error() string {
	return titleNotFoundPrefix + ": " + string(e)
}

type Title map[string]int

func (t Title) names() []string {
	ns := make([]string, len(t))
	for k, n := range t {
		ns[n] = k
	}
	return ns
}

func (t Title) index(names []string) ([]int, error) {
	indexes := make([]int, len(names))
	for i, n := range names {
		if ind, exists := t[n]; exists {
			indexes[i] = ind
		} else {
			return nil, TitleNotFound(fmt.Sprintf("%s cannot be found", n))
		}
	}
	return indexes, nil
}

func (t Title) sortingMarkers(nm []NamedMarker) ([]Marker, error) {
	markers := make([]Marker, len(nm))
	for i, m := range nm {
		if ind, exists := t[m.Name]; exists {
			markers[i] = Marker{ind, m.Order}
		} else {
			return nil, TitleNotFound(fmt.Sprintf("%s cannot be found", m.Name))
		}
	}
	return markers, nil
}

type Processor struct {
	titles Title
	rows   [][]string
}

// Load loads the content of a csv file. The first line of the content will be used
// to define titles.
// Can use os.DirFs(".") as system
// I would like to use os.Open(name string) (*File, error) to interact with csv file in storage
// Files need to be closed after reading. There are things convoluted I have not figured out a way
// to deal with them clearly
// func (p Processor) Load(system fs.FS, name string) {
// 	fmt.Println("Open", name)

// 	f, err := system.Open(name)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer f.Close()

// 	p.Read(f)

// }

// csv.NewReader needs an io.Reader. fs.File defines Reader interface
// os.File is one implementation.
// What I want to start with is a file name.
// To open a file,
func read(source io.Reader) [][]string {
	r := csv.NewReader(source)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}

func (p *Processor) Size() (int, int) {
	if p == nil {
		return 0, 0
	}

	return len(p.titles), len(p.rows)
}

func (p *Processor) Print() {
	fmt.Println("Titles:")
	// Keys returns a slice without a determined order
	fmt.Println(strings.Join(p.titles.names(), ", "))

	fmt.Println("Rows:")
	for i, r := range p.rows {
		fmt.Println(i+1, strings.Join(r, ", "))
	}
	fmt.Println()

}

func (p *Processor) Sort(markers []Marker) {
	sorter := OrderByColumns(markers)
	sorter.Sort(p.rows)
}

func (p *Processor) Export(name string, titles Title, content []string) {

}

// createRecords creates a slice of map with string keys and values
func createRecords(lines [][]string) []map[string]string {
	var records []map[string]string

	for i := 1; i < len(lines); i++ {
		row := make(map[string]string)
		for j, c := range lines[i] {
			row[lines[0][j]] = c
		}
		records = append(records, row)
	}

	return records
}

// createTitle returns a Title from the input of a slice of string.
// The value of each entry is the zero-base column number from a csv file.
func createTitle(names []string) Title {
	t := make(Title)
	for i, n := range names {
		t[n] = i
	}
	return t
}

// NewProcessor opens a csv file named as fileName and returns *Processor
// when there is no error otherwise it logs error and exits
func NewProcessor(fileName string) *Processor {
	fmt.Println("We will be process csv = ", fileName)

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	records := read(file)

	p := &Processor{createTitle(records[0][:]), records[1:][:]}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

	return p
}
