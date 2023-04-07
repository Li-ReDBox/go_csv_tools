package main

import (
	"flag"
	"fmt"
	"os"

	funtool "funmech.com/csv"
)

var csvFile = flag.String("csv", "", "csv file to be processed")

func usage() {
	fmt.Fprintf(os.Stderr, "usage: go run main -csv=some.csv\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	// Check usage.
	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, `Unexpected arguments.`)
		usage()
	}
	if *csvFile == "" {
		usage()
	}

	fmt.Println("We will be process csv = ", *csvFile)

	p := funtool.NewProcessor(*csvFile)

	fmt.Printf("Loaded data from %s, its has size of %v\n", *csvFile, fmt.Sprint(p.Size()))

	p.Print()

	f, err := os.Create("demo.csv")
	if err != nil {
		fmt.Println("Cannot create demo.csv, details:", err)
		os.Exit(1)
	}
	err = p.Write(f)
	if err != nil {
		fmt.Println("Failed to save to a new file: demo.csv. error: ", err)
	}
}
