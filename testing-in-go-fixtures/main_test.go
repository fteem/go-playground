package main

import (
	"os"
	"testing"
)

func TestNewGradebook_ErrorHandling(t *testing.T) {
	cases := []struct {
		fixture   string
		returnErr bool
		name      string
	}{
		{
			fixture:   "testdata/grades/empty.csv",
			returnErr: false,
			name:      "EmptyFile",
		},
		{
			fixture:   "testdata/grades/invalid.csv",
			returnErr: true,
			name:      "InvalidFile",
		},
		{
			fixture:   "testdata/grades/valid.csv",
			returnErr: false,
			name:      "ValidFile",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reader, _ := os.Open(tc.fixture)
			_, err := NewGradebook(reader)
			returnedErr := err != nil

			if returnedErr != tc.returnErr {
				t.Fatalf("Expected returnErr: %v, got: %v", tc.returnErr, returnedErr)
			}
		})
	}
}
