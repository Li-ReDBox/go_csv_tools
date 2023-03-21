package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"funmech.com/csv"
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

	file, err := os.Open(*csvFile)

	if err != nil {
		log.Fatal(err)
	}

	p := csv.Processor{}
	p.Read(file)

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}
