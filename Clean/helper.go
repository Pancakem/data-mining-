package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func loadKeywords(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}

	defer f.Close()

	l := make([]string, 0)
	scanner := bufio.NewScanner(f)

	// read each line of text from the file
	for scanner.Scan() {
		l = append(l, scanner.Text())
	}
	return l
}

// writeCSV writes to CSV file the processed posts
func writeCSV(records [][]string, platform string) {

	// newfilename := changeFileExtension(filename, ".csv")
	newfilename := "processed_" + platform + "_data.csv"
	f, err := os.OpenFile(newfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening write file: ", err)
	}

	defer f.Close()

	writer := csv.NewWriter(f)
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Println(err)
	}
}


// write the items found out from the tweets to file
func writeToFile(filename string, items []string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening write file: ", err)
	}

	for _, v := range items {
		fmt.Fprintln(f, v)
	}

	err = f.Close()
	if err != nil {
		log.Printf("Error writing to %v\n", filename)
	}
}

// get the file extension
// only supported file extensions are .json and .csv
func findFileExtension(filename string) string {
	var extension string

	indx := strings.LastIndex(filename, ".")

	if indx != -1 {
		extension = filename[indx+1:]
	}
	return extension
}

func changeFileExtension(filename string, newExtension string) string {
	var newfilename string
	indx := strings.LastIndex(filename, ".")

	if indx != -1 {
		newfilename = filename[:indx] + newExtension
	}
	return newfilename
}

// removes duplicates from an array
func removeDuplicates(items []string) []string {
	newItems := make([]string, 0)
	keys := make(map[string]bool)

	for _, v := range items {
		v = strings.ToLower(v)
		if _, value := keys[v]; !value {
			keys[v] = true
			newItems = append(newItems, v)
		}
	}
	return newItems
}

func filesInFolder(folderpath string) []string {
	var files []string

	err := filepath.Walk(folderpath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Println("Error reading the folder: ", err)
	}
	// fmt.Println("Success!")
	return files
}

func ifJSON(filename string) bool{
	indx := strings.LastIndex(filename, ".")
	if indx == -1{
		return false
	}

	if filename[indx+1:] == "json"{
		return true
	}

	return false
}

func getDir() string {
	dir, err := os.Getwd()
	if  err != nil {
		log.Fatal(err)
	}

	idx := strings.LastIndex(dir, "/")
	if idx == -1 {
		return ""
	}

	return dir[:idx]
}

func findDataFolder(platform string) string{
	root := getDir()
	datafolder := "/Results/"
	if platform == "twitter"{
		return root + datafolder + "tweet_data/tweet"
	}else if platform == "instagram" {
		return root + datafolder + "instagram_data"
	}
	return ""
}