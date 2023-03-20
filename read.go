package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"log"
	"strings"
)

type Title map[string]int

type Processor struct {
	titles Title
	rows   []string
}

// Load loads the content of a csv file. The first line of the content will be used
// to define titles.
// Can use os.DirFs(".") as system
func (p Processor) Load(system fs.FS, name string) {
	fmt.Println("Open", name)

	f, err := system.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	p.Read(f)

}

// csv.NewReader needs an io.Reader. fs.File defines Reader interface
// os.File is one implementation.
// What I want to start with is a file name.
// To open a file,
func (p Processor) Read(source io.Reader) [][]string {
	r := csv.NewReader(source)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(records)
	return records

}

func (p Processor) Export(name string, titles Title, content []string) {

}

func Read() [][]string {
	in := `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"
`
	r := csv.NewReader(strings.NewReader(in))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(records)
	return records
}

func CreateRecords() []map[string]string {
	var records []map[string]string
	lines := Read()
	row := make(map[string]string)

	for i := 1; i < len(lines); i++ {
		for j, c := range lines[i] {
			row[lines[0][j]] = c
		}
		records = append(records, row)
	}

	return records
}

// CreateTitle returns a Title from the input of a slice of string.
// The value of each entry is the zero-base column number from a csv file.
func CreateTitle(names []string) Title {
	t := make(Title)
	for i, n := range names {
		t[n] = i
	}
	return t
}
