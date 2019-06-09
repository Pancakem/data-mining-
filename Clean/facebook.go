package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type allFbData struct {
	Data []facebook `json:"fields"`
}

type facebook struct {
	ID          string  `json:"id"`
	Message     string  `json:"message"`
	Description string  `json:"description"`
	Comments    comment `json:"comment"`
	Likes       int     `json:"likes"`
}

type comment struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	From    string `json:"from"`
	Likes   int    `json:"like_count"`
}

func loadRecordsFromCSV(filename string) [][]string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}

	defer f.Close()

	records := make([][]string, 0)
	reader := csv.NewReader(f)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		records = append(records, record)
	}

	return records
}
