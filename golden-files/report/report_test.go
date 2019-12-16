package report

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/fteem/go-playground/golden-files/books"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestGenerate(t *testing.T) {
	testcases := []struct {
		name   string
		books  []books.Book
		golden string
	}{
		{
			name: "WithInventory",
			books: []books.Book{
				books.Book{
					Title:     "The Da Vinci Code",
					Author:    "Dan Brown",
					Publisher: "Corgi",
					Pages:     592,
					ISBN:      "978-0552161275",
					Price:     7,
				},
				books.Book{
					Title:     "American on Purpose",
					Author:    "Craig Ferguson",
					Publisher: "Harper Collins",
					Pages:     288,
					ISBN:      "978-0061959158",
					Price:     19,
				},
			},
			golden: "with_inventory",
		},
		{
			name:   "EmptyInventory",
			books:  []books.Book{},
			golden: "empty_inventory",
		},
	}

	for _, testcase := range testcases {
		got := Generate(testcase.books)
		want := goldenValue(t, testcase.golden, got, *update)

		if got != want {
			t.Errorf("Want:\n%s\nGot:\n%s", want, got)
		}
	}
}

func goldenValue(t *testing.T, goldenFile string, actual string, update bool) string {
	t.Helper()
	goldenPath := "testdata/" + goldenFile + ".golden"

	f, err := os.OpenFile(goldenPath, os.O_RDWR, 0644)
	defer f.Close()

	if update {
		_, err := f.WriteString(actual)
		if err != nil {
			t.Fatalf("Error writing to file %s: %s", goldenPath, err)
		}

		return actual
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", goldenPath, err)
	}
	return string(content)
}
