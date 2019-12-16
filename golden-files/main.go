package main

import (
	"fmt"

	"github.com/fteem/go-playground/golden-files/books"
	"github.com/fteem/go-playground/golden-files/report"
)

func main() {
	fmt.Println(report.Generate(books.Books))
}
