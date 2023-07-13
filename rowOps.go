package csv

// replace replaces elements marked by markers by calling operator.
// input: a slice of string; markers: a slice of indexes of elements need to be operated; op: a function: func([]string)
func replace(row []string, markers []int, operator func([]string)) {
	working := make([]string, len(markers))

	for i, m := range markers {
		working[i] = row[m]
	}

	operator(working)

	for i, m := range markers {
		row[m] = working[i]
	}
}

// biop operates on two elements of a slice marked by indxA and indxB, returns a new slice with original elements and one more of the operational result
func biop(row []string, indxA, indxB int, operator func(a, b string) string) []string {
	return append(row, operator(row[indxA], row[indxB]))
}

// Action defines changes to be applied to a slice of string
type Action func(elems []string)

// Operation defines how to check the conditions are met on a slice of string and what to do if the conditions are met.
type Operation struct {
	Check Isfunc
	Act   Action // Act: needs to understand the elements structure.
}

// Do checks the conditions and perform the Act
func (o Operation) Do(elems []string) {
	if o.Check(elems) {
		o.Act(elems)
	}
}
