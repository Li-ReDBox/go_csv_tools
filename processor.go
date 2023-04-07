package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

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

func (p *Processor) Swap(i, j string) error {
	indI, exist := p.titles[i]
	if !exist {
		return TitleNotFound(i)
	}
	indJ, exist := p.titles[j]
	if !exist {
		return TitleNotFound(j)
	}
	for _, r := range p.rows {
		r[indI], r[indJ] = r[indJ], r[indI]
	}
	p.titles[i], p.titles[j] = indJ, indI
	return nil
}

// Extract gets a sub set of current data by a given slice of title names
func (p *Processor) Extract(names []string) ([][]string, error) {
	c := len(names)
	indexes := make([]int, c)
	var ind int
	var exists bool
	for i := 0; i < c; i++ {
		if ind, exists = p.titles[names[i]]; !exists {
			return nil, TitleNotFound(names[i])
		}
		indexes[i] = ind
	}
	rows := len(p.rows)
	extracted := make([][]string, rows)
	for r := 0; r < rows; r++ {
		extracted[r] = make([]string, c)
		for i := 0; i < c; i++ {
			extracted[r][i] = p.rows[r][indexes[i]]
		}
	}
	return extracted, nil
}

// Convert creates a new Processor with a sub set of current data by a given slice of title names
func (p *Processor) Convert(names []string) (*Processor, error) {
	extracted, err := p.Extract(names)
	if err != nil {
		return nil, err
	}
	return &Processor{createTitle(names), extracted}, nil
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
