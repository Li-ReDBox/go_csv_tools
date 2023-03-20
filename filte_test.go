package csv

import (
	"testing"
	"testing/fstest"
)

func TestProcessorLoad(t *testing.T) {
	in := `first_name_file,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","griFile"
`
	fs := fstest.MapFS{
		"authors.csv": {Data: []byte(in)},
	}

	p := Processor{}

	p.Load(fs, "authors.csv")

}
