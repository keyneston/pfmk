package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yuin/goldmark"
)

func main() {
	var inputFile, outputFile string

	flag.StringVar(&inputFile, "i", "", "Input file to be converted, defaults to STDIN")
	flag.StringVar(&outputFile, "o", "", "Where to place output files, defaults to STDOUT")
	flag.Parse()

	var input io.Reader
	var output io.Writer

	if inputFile == "" {
		input = os.Stdin
	} else {
		i, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file: %v", err)
			os.Exit(255)
		}
		input = i
		defer i.Close()
	}
	if outputFile == "" {
		output = os.Stdout
	} else {
		o, err := os.OpenFile(inputFile, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0x644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening output file: %v", err)
			os.Exit(255)
		}
		output = o
		defer o.Close()
	}

	if err := convert(input, output); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing input: %v", err)
		os.Exit(255)
	}
}

func convert(input io.Reader, output io.Writer) error {
	buf := &bytes.Buffer{}

	if _, err := buf.ReadFrom(input); err != nil {
		return fmt.Errorf("Error reading input: %w", err)
	}

	if err := goldmark.Convert(buf.Bytes(), output); err != nil {
		return fmt.Errorf("Error converting input: %w", err)
	}

	return nil
}
