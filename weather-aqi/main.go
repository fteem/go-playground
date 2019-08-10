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

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for group := range groups {
				group.work(rows)
			}
			wg.Done()
		}(&wg)
	}

	go func() {
		wg.Wait()
		close(rows)
	}()

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
