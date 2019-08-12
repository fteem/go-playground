package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type row []string

func readCSV(path string, out chan row) {
	csvFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open CSV: %v", err)
	}

	reader := csv.NewReader(csvFile)

	go func(reader *csv.Reader) {
		for {
			var line, error = reader.Read()
			if error == io.EOF {
				csvFile.Close()
				close(out)
				break
			} else if err != nil {
				log.Fatalf("Error reading file: %v", err)
			}
			out <- line
		}
	}(reader)
}
