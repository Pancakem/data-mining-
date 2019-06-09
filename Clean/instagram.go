package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"strings"
)

type instagramComment struct {
	Author   string   `json:"author"`
	Comment  string   `json:"comment"`
	Mentions []string `json:"mentions"`
	Hashtags []string `json:"hashtags"`
}

type instagram struct {
	ID        string             `json:"key"`
	Likes     int                `json:"likes"`
	Caption   string             `json:"caption"`
	Comments  []instagramComment `json:"comments"`
	Mentions  []string           `json:"mentions"`
	Hashtags  []string           `json:"hashtags"`
	Timestamp string             `json:"datetime"`
}

func loadInstagramData(files []string) ([][]string, []string, []string) {
	var (
		records  [][]string
		mentions []string
		hashtags []string
	)
	ips := readFiles(files)

	for k, v := range ips {
		recrds, mntions, hshtags := toArray(v, files[k])

		for _, v := range recrds {
			records = append(records, v)
		}

		for _, v := range mntions {
			mentions = append(mentions, v)
		}

		for _, v := range hshtags {
			hashtags = append(hashtags, v)
		}
	}

	return records, mentions, hashtags
}

func getNameFromFile(filename string) string {
	return filename[:strings.LastIndex(filename, ".")]
}

func readFiles(files []string) [][]*instagram {
	instaPosts := make([][]*instagram, 0)

	for _, file := range files {
		var idata []*instagram

		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(data, &idata)
		if err != nil {
			log.Fatal(err)
		}

		instaPosts = append(instaPosts, idata)

	}
	return instaPosts
}

func toArray(ip []*instagram, filename string) (records [][]string, hashtags []string, mentions []string) {

	for _, v := range ip {
		record := []string{v.ID, getNameFromFile(filename), v.Caption, strings.Join(v.Mentions, " "), strconv.Itoa(v.Likes), v.Timestamp}
		records = append(records, record)

		for _, h := range v.Hashtags {
			hashtags = append(hashtags, h)
		}

		for _, m := range v.Mentions {
			mentions = append(mentions, m)
		}

		// get all comments
		for _, c := range v.Comments {
			record = []string{v.ID, c.Author, c.Comment, strings.Join(c.Mentions, " "), "0", v.Timestamp}
			records = append(records, record)
			for _, m := range c.Mentions {
				mentions = append(mentions, m)
			}
			for _, h := range c.Hashtags {
				hashtags = append(hashtags, h)
			}
		}

	}

	return records, hashtags, mentions
}
