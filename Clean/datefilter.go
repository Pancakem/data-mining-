package main

import (
	"log"
	"strings"
	"time"
)

// in format from twitter `2008-01-24 22:39:47`
// in format from instagram `2019-05-12T17:29:54.000Z`- no need to format
func formatTimeStamp(timestamp string) string {
	//append .000Z to the end
	timestamp += ".000Z"
	newString := strings.Replace(timestamp, " ", "T", 1)
	return newString
}

func getTimeFromStamp(timestamp string) *time.Time {
	timeValue, err := time.Parse("2019-05-12T17:29:54.000Z", timestamp)
	if err != nil {
		log.Fatal(err)
	}
	return &timeValue
}

// isInLimit finds if a timestamp adheres to a set time period
// limits is a 2 element slice with the indices being the limits
func isInLimit(limits []string, timestamp string) bool {
	upperLimit := getTimeFromStamp(limits[1]).Unix()
	lowerLimit := getTimeFromStamp(limits[0]).Unix()

	theTime := getTimeFromStamp(timestamp).Unix()

	// if the timestamp is within the limits
	return theTime >= lowerLimit && theTime <= upperLimit
}
