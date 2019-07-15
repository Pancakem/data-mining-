package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
)

func main() {
	var (
		offset           int
		processedRecords [][]string
		users            []string
		platform         string
		folderpath       string
		hashtags         []string
		mentions         []string
	)
	var wg sync.WaitGroup

	flag.StringVar(&platform, "platform", "", "enter the platform being processed")

	flag.Parse()

	folderpath = findDataFolder(platform)
	tempfiles := filesInFolder(folderpath)

	files := make([]string, len(tempfiles))
	copy(files, tempfiles)
	offset = 2
	
	

	keywords := loadKeywords("search_words.txt")

	jobs, mentions, hashtags := splitIntoJobs(files, platform)

	rtChan := make(chan *job, len(jobs))
	ctChan := make(chan int)
	wg.Add(1)
	go func() {
		pr := merge(rtChan, ctChan, len(jobs))

		for _, v := range pr {
			processedRecords = append(processedRecords, v)
		}

		wg.Done()
	}()

	g := len(jobs)

	wg.Add(g)

	for _, jb := range jobs {

		go func(gojob *job) {

			defer wg.Done()

			j := work(keywords, gojob, offset)
			if len(j.records) > 0 {
				rtChan <- j
			}
			if platform == "twitter" {
				var userss []string
				for _, record := range j.records {
					userss = findUsers(record[offset])

					if len(userss) != 0 {
						for _, v := range userss {
							users = append(users, v)
						}
					}
				}
			}
			ctChan <- 1

		}(jb)
	}

	wg.Wait()
	if platform == "twitter" {
		// remove duplicates
		newRecords := removeDuplicateTweets(processedRecords)
		sortedRecords := quickSort(newRecords)
		// newUser := removeDuplicates(users)
		err := writeCSV(sortedRecords, "data.csv")
		if err != nil {
			log.Println("Error occured in xlsx file", err)
		}

		// writeToFile("users.txt", newUser)

	} else if platform == "instagram" {
		root := getDir()
		newHashtags := removeDuplicates(hashtags)
		writeToFile(root+"/Parameters/hashtags.txt", newHashtags)
		newMentions := removeDuplicates(mentions)
		writeToFile(root+"/Parameters/mentions.txt", newMentions)
		newRecords := removeDuplicateTweets(processedRecords)
		err := writeCSV(newRecords, "temp.csv")
		if err != nil {
			log.Println("Couldn't write csv file")
		}

	}

	err := generateXLSXFromCSV("temp.csv", "data.xlsx", ",")
	if err != nil {
		log.Println("Couldn't create xlsx file")
	}
	//deleteFile("data.csv")

	fmt.Println("Processing complete.")
}
