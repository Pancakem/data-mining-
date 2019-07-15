package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/tealeg/xlsx"
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
	//fmt.Println("Success!")
	//fmt.Println(len(files))
	return files
}

func ifJSON(filename string) bool {
	indx := strings.LastIndex(filename, ".")
	if indx == -1 {
		return false
	}

	if filename[indx+1:] == "json" {
		return true
	}

	return false
}

func getDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//idx := strings.LastIndex(dir, "/")
	//if idx == -1 {
	//	return ""
	//}

	//return dir[:idx]
	return dir
}

func findDataFolder(platform string) string {
	root := getDir()
	datafolder := "\\Data\\tweet"
	if platform == "twitter" {
		return root + datafolder //+ "tweet_data/tweet"
	} else if platform == "instagram" {
		return root + datafolder + "instagram_data"
	}
	return ""
}

// writeCSV writes to CSV file the processed posts
func writeCSV(records [][]string, filename string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening write file: ", err)
		return err
	}

	// read cache and add to the to be written records
	files := filesInFolder(cacheFolder())

	for _, file := range files[1:] {
		records = append(records, readCacheFile(file))
	}

	defer f.Close()

	writer := csv.NewWriter(f)
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			cacheWriteFailure(record)
			log.Fatalln("error writing record to csv:", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Println(err)
	}
	return err
}

// cacheWriteFailure takes the row
// provides caching for failed  write
// the cache is a file of the row data
func cacheWriteFailure(record []string) {
	// provide a random uuid4 string for each file
	u2, err := uuid.NewV4()
	if err != nil {
		log.Println("Error generating uuid: ", err)
	}

	uid := u2.String()

	f, err := os.OpenFile(cacheFolder()+uid, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Couldn't create cache file")
	}
	for _, v := range record {
		fmt.Fprintln(f, v)
	}

	defer f.Close()
}

func cacheFolder() string {
	return getDir() + "/Cache/"
}

func readCacheFile(file string) []string {
	record := make([]string, 0)

	f, err := os.Open(file)
	if err != nil {
		log.Println("Couldn't not open file", err)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		record = append(record, scanner.Text())
	}
	return record
}

func generateXLSXFromCSV(csvPath string, XLSXPath string, delimiter string) error {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	if len(delimiter) > 0 {
		reader.Comma = rune(delimiter[0])
	} else {
		reader.Comma = rune(';')
	}
	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet(csvPath)
	if err != nil {
		return err
	}
	fields, err := reader.Read()
	for err == nil {
		row := sheet.AddRow()
		for _, field := range fields {
			cell := row.AddCell()
			cell.Value = field
		}
		fields, err = reader.Read()
	}
	if err != nil {
		fmt.Printf(err.Error())
	}
	return xlsxFile.Save(XLSXPath)
}

func deleteFile(file string) error {
	return os.Remove(file)
}
