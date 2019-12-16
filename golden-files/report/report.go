package report

import (
	"bytes"
	"log"
	"text/template"

	"github.com/fteem/go-playground/golden-files/books"
)

const (
	header string = `
| Title         | Author        |  Pages  |  ISBN  |  Price  |
| ------------- | ------------- | ------- | ------ | ------- |
`
	rowTemplate string = "| {{ .Title }} | {{ .Author }} | {{ .Pages }} | {{ .ISBN }} | {{ .Price }} |"
)

func Generate(books []books.Book) string {
	buf := bytes.NewBufferString(header)

	t := template.Must(template.New("table").Parse(rowTemplate + "\n"))

	for _, book := range books {
		err := t.Execute(buf, book)
		if err != nil {
			log.Println("Error executing template:", err)
		}
	}
	return buf.String()
}
