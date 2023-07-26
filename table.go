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

// NewTable opens a csv file named by fileName and returns *Table
// when there is no error, otherwise it logs error and exits
func NewTable(fileName string) *Table {
	fmt.Println("We will be process csv = ", fileName)

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	records, err := read(file)
	ce := file.Close()

	if err != nil {
		log.Fatal(err)
	}

	if ce != nil {
		log.Fatal(ce)
	}

	return &Table{createTitle(records[0][:]), records[1:][:]}
}

// Table is a data structure, so the name is not a good choice
// The data are rows, with (optional) titles, more like a table.
type Table struct {
	// titles is optional internally, but as column methods needs them to work correctly,
	// the NewTable function automatically take the first row as the titles.
	titles Title
	rows   [][]string
}

// read is a wrapper of csv.Reader.ReadAll.
// csv.NewReader needs an io.Reader. fs.File defines Reader interface
// os.File is one implementation and strings.NewReader is another one.
// It exists for supporting tests.
func read(source io.Reader) ([][]string, error) {
	r := csv.NewReader(source)

	return r.ReadAll()
}

// Size returns the number of columns and the number of rows of matrix rows in that order.
func (p *Table) Size() (int, int) {
	if p == nil {
		return 0, 0
	}

	return len(p.rows[0]), len(p.rows)
}

// Print prints the data in two sections: Titles, if there is no titles, a blank line; and Rows.
func (p *Table) Print() {
	fmt.Println("Titles:")
	// Keys returns a slice without a determined order
	fmt.Println(strings.Join(p.titles.names(), ", "))

	fmt.Println("Rows:")
	for i, r := range p.rows {
		fmt.Println(i+1, strings.Join(r, ", "))
	}
	fmt.Println()

}

// Sort sorts the rows according to Markers: which column, in what direction.
func (p *Table) Sort(markers []Marker) {
	sorter := OrderByColumns(markers)
	sorter.Sort(p.rows)
}

// Swap swaps two columns identified by their names. If any of name is not found, TitleNotFound error returns.
// Swap definitely needs titles.
func (p *Table) Swap(i, j string) error {
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

// Extract extracts columns named by given title names, the order of extracted columns
// is defined by the order of given names argument. This can be used to reorder columns
// when lengths of names equals to titles.
// The returned [][]string is independent to its source.
// Extract definitely needs titles.
func (p *Table) Extract(names []string) ([][]string, error) {
	inds, err := p.titles.indexes(names)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Extract method: %w", err)
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

// Convert creates a new Table by extracting named columns. The order of resulting columns is
// defined by the order of given names argument. This can be used to reorder columns
// when lengths of names equals to titles. The returned Table is independent to its source.
// Convert definitely needs titles.
func (p *Table) Convert(names []string) (*Table, error) {
	extracted, err := p.Extract(names)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Convert method: %w", err)
	}
	return &Table{createTitle(names), extracted}, nil
}

// Split uses the values of columns identified by title names to group rows and creates a slice of new Tables.
// The source Table should have been sorted, the order of names is significant.
// The new Tables are linked to the source, so they are views of original.
// Split definitely needs titles.
func (p *Table) Split(names []string) ([]*Table, error) {
	inds, err := p.titles.indexes(names)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Split method: %w", err)
	}

	var np []*Table
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
			// any checker is different, it means a new Table
			if p.rows[r][inds[i]] != current[i] {
				// slice a block of rows to create a new Table and append to the returning slice.
				np = append(np, &Table{p.titles, p.rows[start:r]})
				update(r)
				start = r
				break
			}
		}
	}
	if start < len(p.rows) {
		np = append(np, &Table{p.titles, p.rows[start:]})
	}

	return np, nil
}

// Write the data to the Writer w
func (p *Table) Write(w io.Writer) error {
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

type Isfunc func(elems []string) bool

// Filter uses condition functions to check rows and remove them if all conditions are met.
// This a in place procedure: p.rows are replaced.
func (p *Table) Filter(is []Isfunc) {
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

// Unique creates a new Table from the current one by removing all the duplicates
func (p *Table) Unique() *Table {
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

	return &Table{p.titles, unique}
}

// Clone makes a complete new Table from the current one, so both can be processed independently.
func (p *Table) Clone() *Table {
	nCols, nRows := p.Size()
	r := make([][]string, 0, nRows)
	for i := 0; i < nRows; i++ {
		c := make([]string, nCols)
		copy(c, p.rows[i])
		r = append(r, c)
	}
	return &Table{p.titles.clone(), r}
}

// createRecords creates a slice of map by turning each line from the second line onwards into a map with string keys come from the first line.
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

// Replace replaces some elements of rows. Each operation has a condition checker and an act to be performed when the conditions are met.
func (p *Table) Replace(ops []Operation) {
	nCols, nRows := p.Size()
	if nCols == 0 || nRows == 0 {
		fmt.Printf("Cannot operate on a zero sized data set: nCols = %d, nRows = %d\n", nCols, nRows)
		return
	}

	for i := 0; i < nRows; i++ {
		for _, op := range ops {
			op.Do(p.rows[i])
		}
	}
}

// Derive use two existing columns to derive and add the result to a new column to the end
// Derive definitely needs titles
func (p *Table) Derive(inxA, inxB int, name string, op func(a, b string) string) {
	if _, ok := p.titles[name]; ok {
		fmt.Println(name, ": cannot existing column name as a new colum")
		return
	}
	p.titles[name] = len(p.titles)

	for i := 0; i < len(p.rows); i++ {
		p.rows[i] = append(p.rows[i], op(p.rows[i][inxA], p.rows[i][inxB]))
	}
}
