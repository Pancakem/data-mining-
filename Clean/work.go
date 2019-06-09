package main

import (
	"fmt"
	"strings"
)

type job struct {
	records [][]string
}

func newjob() *job {
	return &job{}
}

// splits the records into various jobs to be run concurrently
func splitIntoJobs(files []string, social string) ([]*job, []string, []string) {
	var (
		records  [][]string
		mentions []string
		hashtags []string
	)

	if social == "twitter" {
		records = loadAllTweets(files)
	} else if social == "instagram" {
		records, mentions, hashtags = loadInstagramData(files)
	}

	numofJobs := len(records)
	jobs := make([]*job, 0)
	splitNum := int(numofJobs / 100)
	rem := numofJobs % 100
	if splitNum < 100 {
		j := job{records: records}
		jobs = append(jobs, &j)
		return jobs, mentions, hashtags
	}
	j := job{records: records[:(splitNum + rem)]}
	records = records[(splitNum + rem):]
	jobs = append(jobs, &j)
	for i := 1; i < 100; i++ {
		j := job{records: records[:splitNum]}
		records = records[splitNum:]
		jobs = append(jobs, &j)
	}
	return jobs, mentions, hashtags
}

// get a new filename for the processed data
// append `processed_` before the filename
func newFileName(filename string) string {
	return "processed_" + filename
}

// merge puts together the jobs that have passed the check
func merge(rt chan *job, ct chan int, jobLength int) [][]string {
	records := make([][]string, 0)
	count := 0
	total := 0

	for total < jobLength {
		select {
		case x := <-rt:
			count++
			for _, v := range x.records {
				records = append(records, v)
			}
		case y := <-ct:
			total += y
		}

	}
	fmt.Println("Goroutines done: ", total)
	fmt.Printf("Channel read in %v jobs\n", count)

	return records
}

// work matches bully keyword in the text
func work(keywords []string, j *job, offset int) *job {
	newJob := newjob()

	for _, v := range j.records {
		// get the relevant text string from the array
		// for simplicity we'll use some offset number
		// to tell the position of the text in the array
		text := v[offset]

		if detectEnglish(text) {
			for _, kw := range keywords {
				if strings.Contains(strings.ToLower(text), kw) {
					newJob.records = append(newJob.records, v)
					break
				}
			}
		}
	}
	return newJob
}
