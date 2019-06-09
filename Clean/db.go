package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

func writeDB(records [][]string) error {
	db := initDB()

	query := "INSERT INTO posts (id, username, post_text, towhom, likes, time_stamp) VALUES ($1,$2,$3,$4,$5,$6)"
	for _, record := range records {
		writedb(record, db, query)
	}
	return nil
}

func writedb(record []string, db *sql.DB, query string) {
	err := writeDBRow(db, record, query)
	if err != nil {
		log.Println("Error writing to db caching")
		cacheDBWriteFailure(record)
	}
}

func initDB() *sql.DB {
	connectStr := "user=postgres dbname=postgres password=Apassword host=localhost port=5432 sslmode=disable"
	db, _ := sql.Open("postgres", connectStr)
	err := db.Ping()
	if err != nil {
		log.Fatal("Could not create database connection: ", err)
	}

	return db
}

func writeDBRow(db *sql.DB, record []string, query string) error {

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Could not prepare the sql statement:", err)
		cacheDBWriteFailure(record)
		return err
	}

	_, err = stmt.Exec(record[0], record[1], record[2], record[3], record[4], record[5])
	if err != nil {
		log.Println("Could not execute the sql statement:", err)
		cacheDBWriteFailure(record)
		return err
	}

	return nil
}

// cacheDBWriteFailure takes the row
// provides caching for failed db write
// the cache is a file of the row data
func cacheDBWriteFailure(record []string) {
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

func retryCache(db *sql.DB) error {

	// walk the cache folder for rows
	caches := filesInFolder(cacheFolder())

	for _, cache := range caches {
		record := readCacheFile(cache)
		err := writeDBRow(db, record, "") // write query here
		if err != nil {
			log.Println("Failed to redo cache", err)
			continue
		}
		err = deleteFile(cache)
		if err != nil {
			log.Println("Could not delete cache file:", err)
		}
	}

	return nil
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

func deleteFile(file string) error {
	return os.Remove(file)
}
