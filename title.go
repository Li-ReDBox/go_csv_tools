package csv

import (
	"fmt"
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

func (t Title) indexes(names []string) ([]int, error) {
	inds := make([]int, len(names))
	for i, n := range names {
		if ind, exists := t[n]; exists {
			inds[i] = ind
		} else {
			return nil, TitleNotFound(fmt.Sprintf("%s cannot be found", n))
		}
	}
	return inds, nil
}

// sortingMarkers creates a slice of Marker from a give slice of NamedMarker
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

// createTitle returns a Title from the input of a slice of string.
// The value of each entry is the zero-base column number from a csv file.
func createTitle(names []string) Title {
	t := make(Title)
	for i, n := range names {
		t[n] = i
	}
	return t
}
