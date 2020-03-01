package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Snippet struct {
	gorm.Model

	Body     string `gorm:"type:text"`
	URI      string `gorm:"type:varchar(255);unique"`
	Title    string
	Language string
}

func (s *Snippet) BeforeCreate() (err error) {
	letters := []string{"s", "n", "i", "p", "p", "e", "t"}
	timestamp := time.Now().UnixNano()
	s.URI = fmt.Sprintf("%s-%d", letters[int(timestamp)%len(letters)], timestamp)
	return
}
