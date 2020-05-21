package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ledongthuc/pdf"
)

func persist(content []byte, w io.Writer) error {
	_, err := w.Write(content)
	if err != nil {
		log.Println("Failed to persist contents.")
		return err
	}

	return nil
}

func slurp(r *pdf.Reader) []byte {
	var bs []byte

	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				bs = append(bs, []byte(word.S)...)
			}
		}
	}
	return bs
}

func run(args []string, out io.Writer) error {
	log.SetOutput(out)

	if len(args) < 3 {
		return fmt.Errorf("Expected at least 2 arguments, got %d.", len(args)-1)
	}

	pdfFile, r, err := pdf.Open(args[1])
	if err != nil {
		return err
	}
	defer pdfFile.Close()

	contents := slurp(r)

	txtFile, err := os.Create(args[2])
	if err != nil {
		return err
	}
	defer txtFile.Close()

	err = persist(contents, txtFile)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
