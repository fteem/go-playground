package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ledongthuc/pdf"
)

func TestPersist(t *testing.T) {
	tt := []struct {
		name    string
		content []byte
		out     func() (io.ReadWriter, error)
	}{
		{
			name:    "WithNoContent",
			content: []byte{},
			out: func() (io.ReadWriter, error) {
				return os.Create(filepath.Join(t.TempDir(), "empty.txt"))
			},
		},
		{
			name:    "WithContent",
			content: []byte{},
			out: func() (io.ReadWriter, error) {
				return os.Create(filepath.Join(t.TempDir(), "not-empty.txt"))
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			f, err := tc.out()
			if err != nil {
				t.Fatalf("Cannot create output file: %s", err)
			}

			err = persist(tc.content, f)
			if err != nil {
				t.Fatalf("Cannot persits to output file: %s", err)
			}

			b := []byte{}
			if _, err = io.ReadFull(f, b); err != nil {
				t.Fatalf("Cannot read test output file: %s", err)
			}

			if !bytes.Equal(b, tc.content) {
				t.Errorf("Persisted content is different than saved content.")
			}
		})
	}
}

func TestSlurp(t *testing.T) {
	tt := []struct {
		name    string
		pdfPath string
		size    int
	}{
		{
			name:    "PDFWithContent",
			pdfPath: "testdata/content.pdf",
			size:    11463,
		},
		{
			name:    "PDFWithoutContent",
			pdfPath: "testdata/empty.pdf",
			size:    0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pdfFile, r, err := pdf.Open(tc.pdfPath)
			if err != nil {
				t.Fatalf("Couldn't open PDF %s, error: %s", tc.pdfPath, err)
			}
			defer pdfFile.Close()

			contents := slurp(r)

			if len(contents) != tc.size {
				t.Errorf("Expected contents to be %d bytes, got %d", tc.size, len(contents))
			}
		})
	}
}

func TestRun(t *testing.T) {
	tt := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "WithValidArguments",
			input:  "testdata/input.pdf",
			output: filepath.Join(t.TempDir(), "output.txt"),
		},
		{
			name:   "WithEmptyInput",
			input:  "testdata/empty.pdf",
			output: filepath.Join(t.TempDir(), "output.txt"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := run([]string{"foo", tc.input, tc.output}, os.Stdout)
			if err != nil {
				t.Fatalf("Expected no error, got:  %s", err)
			}

			if _, err := os.Stat(tc.output); os.IsNotExist(err) {
				t.Errorf("Expected persisted file at %s, did not find it: %s", tc.output, err)
			}
		})
	}
}

func TestRunErrors(t *testing.T) {
	tt := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "WithoutArguments",
			input:  "",
			output: "",
		},
		{
			name:   "WithoutOneArgument",
			input:  "testdata/input.pdf",
			output: "",
		},
		{
			name:   "WithNonexistentInput",
			input:  "testdata/nonexistent.pdf",
			output: "testdata/output.txt",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := run([]string{"foo", tc.input, tc.output}, os.Stdout)

			if err == nil {
				t.Fatalf("Expected an error, did not get one.")
			}
		})
	}
}
