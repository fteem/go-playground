package main

import (
	"encoding/csv"
	"log"
	"os"
	"sync"
)

func processLines(in chan row, out chan workerGroup) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for line := range in {
				out <- newWorkerGroup(
					line[2],
					newAqiClient(line[2]),
					newDarkskyClient(line[0], line[1]),
				)
			}
			wg.Done()
		}(&wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}

func processGroups(in chan workerGroup, out chan row) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for group := range in {
				group.work(out)
			}
			wg.Done()
		}(&wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}

func main() {
	lines := make(chan row)
	groups := make(chan workerGroup)
	rows := make(chan row)

	processLines(lines, groups)

	processGroups(groups, rows)

	readCSV("cities.csv", lines)

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
