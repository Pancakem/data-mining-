package main

import (
	"log"
	"math/rand"
	"strings"
	"time"
)

// in format from twitter `2008-01-24 22:39:47`
// in format from instagram `2019-05-12T17:29:54.000Z`- no need to format
func formatTimeStamp(timestamp string) string {
	//append .000Z to the end
	timestamp += "Z"
	newString := strings.Replace(timestamp, " ", "T", 1)
	return newString
}

func getTimeFromStamp(timestamp string) *time.Time {
	timeValue, err := time.Parse(time.RFC3339, timestamp)
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

// create an order in the tweets
func createOrder(records [][]string) [][]string {
	return records
}

// quicksort the timestamps into and ascending order
func quickSort(slice [][]string) [][]string {

	timestampIndex := 5
	length := len(slice)

	if length <= 1 {
		sliceCopy := make([][]string, length)
		copy(sliceCopy, slice)
		return sliceCopy
	}

	m := slice[rand.Intn(length)]

	less := make([][]string, 0, length)
	middle := make([][]string, 0, length)
	more := make([][]string, 0, length)

	for _, item := range slice {
		switch {
		case topCompare(item[timestampIndex], m[timestampIndex], "lesser"): //item < m:
			less = append(less, item)
		case topCompare(item[timestampIndex], m[timestampIndex], "equal"): //item < m:
			middle = append(middle, item)
		case topCompare(item[timestampIndex], m[timestampIndex], "greater"): // item > m:
			more = append(more, item)
		}
	}

	less, more = quickSort(less), quickSort(more)

	less = append(less, middle...)
	less = append(less, more...)

	return less
}

// compares the times
// if time1 is greater it returns a true
// false otherwise
func compare(time1, time2 *time.Time, ty string) bool {
	switch ty {
	case "equal":
		return time1.Unix() == time2.Unix()
	case "greater":
		return time1.Unix() > time2.Unix()
	case "lesser":
		return time1.Unix() < time2.Unix()
	default:
		return false
	}
}

// wrap compare
func topCompare(timestamp1, timestamp2 string, ty string) bool {
	time1 := getTimeFromStamp(formatTimeStamp(timestamp1))
	time2 := getTimeFromStamp(formatTimeStamp(timestamp2))
	return compare(time1, time2, ty)
}
