package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type otherTweet struct {
	Username         string `json:"usernameTweet"`
	ID               string `json:"ID"`
	Text             string `json:"text"`
	Isreply          bool   `json:"is_reply"`
	IsRetweet        bool   `json:"is_retweet"`
	URL              string `json:"url"`
	NumberOfReply    int    `json:"nbr_reply"`
	NumberOfRetweets int    `json:"nbr_retweet"`
	NumberOfFave     int    `json:"nbr_favorite"`
	Time             string `json:"datetime"`
	UserID           string `json:"user_id"`
}


func loadAllTweets(files []string) [][]string {
	tweets := make([]*otherTweet, 0)
	for _, file := range files[2:] {
		tweets = append(tweets, loadJSONOtherTweet(file))
	}

	return toArrayofArrays(tweets)
}

func loadJSONOtherTweet(filename string) *otherTweet {
	var ot otherTweet

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error in reading json file: ", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Error loading file into buffer: ", err)
	}
	defer f.Close()
	
	err = json.Unmarshal(data, &ot)
	if err == nil {
		//log.Println("Error reading buffer: ", err)
		deleteFile(filename)
	}
	return &ot
}

func toArrayofArrays(ots []*otherTweet) (rec [][]string) {

	for _, v := range ots {
		allToWhom := strings.Join(findUsers(v.Text), " ")
		rec = append(rec, []string{v.ID, v.Username, v.Text,
			allToWhom, strconv.Itoa(v.NumberOfRetweets + v.NumberOfFave),
			v.Time,
		})
	}
	return rec
}

func removeDuplicateTweets(items [][]string) [][]string {
	newItems1 := make([][]string, 0)
	encountered := make(map[string][]string)

	// removing duplicates by tweetId
	for _, v := range items {
		str := string(v[0])
		encountered[str] = v
	}

	for _, v := range encountered {
		newItems1 = append(newItems1, v)
	}
	encountered = make(map[string][]string)

	// removing duplicates by post
	newItems := make([][]string, 0)
	for _, v := range newItems1 {
		str := string(v[2])
		encountered[str] = v
	}

	for _, v := range encountered {
		newItems = append(newItems, v)
	}

	return newItems
}

func findUsers(text string) []string {
	var users []string

	at := "@"

	if strings.Contains(text, at) {
		words := strings.Split(text, " ")

		for _, word := range words {
			if strings.Contains(word, at) {
				users = append(users, word)
			}
		}
	}

	return users
}
