package csv

import (
	"crypto/md5"
	"encoding/csv"
	"errors"
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

	n := len(p.titles)

	// If titles is not set, get it from the first row, assume rows are equal
	if n == 0 {
		n = len(p.rows[0])
	}

	return n, len(p.rows)
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

// Extract gets a sub set of current data by a given slice of title names, the order of columns can be changed.
// FIXME: new Processor is linked to the source
func (p *Processor) Extract(names []string) ([][]string, error) {
	inds, err := p.titles.indexes(names)
	if err != nil {
		return nil, err
	}

	c := len(names)
	rows := len(p.rows)
	extracted := make([][]string, rows)

	for r := 0; r < rows; r++ {
		extracted[r] = make([]string, c)
		for i := 0; i < c; i++ {
			extracted[r][i] = p.rows[r][inds[i]]
		}
	}
	return extracted, nil
}

// Convert creates a new Processor with a sub set of current data by a given slice of title names. The order of columns can be changed.
// The new Processor is linked to the source, so it is a view.
func (p *Processor) Convert(names []string) (*Processor, error) {
	extracted, err := p.Extract(names)
	if err != nil {
		return nil, err
	}
	return &Processor{createTitle(names), extracted}, nil
}

// Split uses title names to group rows and creates a slice of new Processor
// The source Processor should have been sorted, the order of names is significant.
// The new Processor is linked to the source, so it is a view.
func (p *Processor) Split(names []string) ([]*Processor, error) {
	inds, err := p.titles.indexes(names)
	if err != nil {
		return nil, err
	}

	var np []*Processor
	c := len(names)
	current := make([]string, c)

	update := func(r int) {
		for i, j := range inds {
			current[i] = p.rows[r][j]
		}
	}

	update(0)
	start := 0
	crows := len(p.rows)

	for r := 1; r < crows; r++ {
		for i := 0; i < c; i++ {
			// any checker is different, it means a new Processor
			if p.rows[r][inds[i]] != current[i] {
				np = append(np, &Processor{p.titles, p.rows[start:r]})
				update(r)
				start = r
				break
			}
		}
	}
	if start < len(p.rows) {
		np = append(np, &Processor{p.titles, p.rows[start:]})
	}

	return np, nil
}

// Write the data to the Writer w
func (p *Processor) Write(w io.Writer) error {
	var tErr, lErr, fErr error
	writer := csv.NewWriter(w)
	names := p.titles.names()
	if len(names) > 0 {
		tErr = writer.Write(names)
	}

	if tErr == nil {
		lErr = writer.WriteAll(p.rows)

		fErr = writer.Error()
	}
	return errors.Join(tErr, lErr, fErr)
}

type Isfunc func(s []string) bool

// Remove uses condition functions to check rows and remove them if conditions are met.
// This a in place procedure.
func (p *Processor) Remove(is ...Isfunc) {
	var (
		i   int
		can bool
	)

	temp := p.rows[:0]
	s := len(is)
	for r := 0; r < len(p.rows); r++ {
		can = false
		for i = 0; i < s; i++ {
			if is[i](p.rows[r]) {
				can = true
				break
			}
		}
		if can {
			temp = append(temp, p.rows[r])
		}
	}
	p.rows = temp
}

func md5hash(o []string) string {
	sum := md5.Sum([]byte(strings.Join(o, "")))
	return fmt.Sprintf("%x", sum)
}

// Unique creates a new Processor from the current one by removing all the duplicates
func (p *Processor) Unique() *Processor {
	markers := make(map[string]struct{})

	nCols, nRows := p.Size()
	unique := make([][]string, 0, nRows)

	for i := 0; i < nRows; i++ {
		h := md5hash(p.rows[i])
		if _, ok := markers[h]; !ok {
			c := make([]string, nCols)
			copy(c, p.rows[i])
			unique = append(unique, c)
			markers[h] = struct{}{}
		}
	}

	return &Processor{p.titles, unique}
}

// Clone makes a complete new Processor from the current one, so both can be processed independently.
func (p *Processor) Clone() *Processor {
	nCols, nRows := p.Size()
	r := make([][]string, 0, nRows)
	for i := 0; i < nRows; i++ {
		c := make([]string, nCols)
		copy(c, p.rows[i])
		r = append(r, c)
	}
	return &Processor{p.titles.clone(), r}
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
