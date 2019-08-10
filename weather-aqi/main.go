package main

import (
	"encoding/csv"
	"log"
	"os"
	"sync"
)

type row []string

func main() {
	groups := make(chan workerGroup)
	rows := make(chan row)

	var wg sync.WaitGroup

	go func(wg *sync.WaitGroup) {
		for group := range groups {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				group.work(rows)
				wg.Done()
			}(wg)
		}
		wg.Wait()
		close(rows)
	}(&wg)

	readCSV("cities.csv", groups)

	outCSV, err := os.Create("result.csv")
	if err != nil {
		log.Fatalf("Cannot create CSV: %v", err)
	}
	defer outCSV.Close()

	writer := csv.NewWriter(outCSV)
	defer writer.Flush()

	for row := range rows {
		err := writer.Write(row)
		if err != nil {
			log.Fatalf("Cannot write to file: %s", err)
		}
	}
}
