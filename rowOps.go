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
