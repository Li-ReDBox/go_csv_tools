package csv

import (
	"testing"

	"golang.org/x/exp/maps"
)

func TestCreateTitle(t *testing.T) {
	names := [...]string{"a", "b", "c"}

	want := Title{"a": 0, "b": 1, "c": 2}

	titles := createTitle(names[:])

	if !maps.Equal(titles, want) {
		t.Errorf("CreateTitle() = %v, want %v", titles, want)
	}
}

func TestTitleclone(t *testing.T) {
	source := createTitle([]string{"a", "b", "c"})
	cloned := source.clone()

	if &source == &cloned {
		t.Errorf("The address are the same: %p == %p\n", &source, &cloned)
	}

	if !maps.Equal(source, cloned) {
		t.Errorf("Title.clone filed: want %v, got = %v\n", source, cloned)
	}
}
