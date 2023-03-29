package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type Title map[string]int

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
	for t := range p.titles {
		fmt.Printf("%s\t", t)
	}
	fmt.Println()

	fmt.Println("Rows:")
	for i, r := range p.rows {
		fmt.Printf("%d, %v\n", i+1, r)
	}
	fmt.Println()

}

func (p *Processor) Sort(markers []int) {
	sorter := OrderByColumns()
	sorter.Sort(p.rows, markers)
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
