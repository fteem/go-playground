package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// Record is a gradebook record
type Record struct {
	student string
	subject string
	grade   string
}

// Gradebook is a collection of Records
type Gradebook []Record

// NewGradebook is a constructor for Gradebook
func NewGradebook(csvFile io.Reader) (Gradebook, error) {
	var gradebook Gradebook
	reader := csv.NewReader(csvFile)

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return gradebook, err
		}

		if len(line) < 3 {
			return gradebook, fmt.Errorf("Invalid file structure")
		}

		gradebook = append(gradebook, Record{
			student: line[0],
			subject: line[1],
			grade:   line[2],
		})
	}

	return gradebook, nil
}

// FindByStudent finds students by their name
func (gb *Gradebook) FindByStudent(student string) []Record {
	var records []Record
	for _, record := range *gb {
		if student == record.student {
			records = append(records, record)
		}
	}
	return records
}

func main() {
	csvFile, err := os.Open("grades.csv")
	if err != nil {
		fmt.Println(fmt.Errorf("error opening file: %v", err))
	}
	grades, err := NewGradebook(csvFile)
	fmt.Printf("%+v\n", grades.FindByStudent("Jane"))
}
