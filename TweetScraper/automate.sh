#!/bin/bash

file="queries.txt"
while read line; do
    echo "\n"
    echo "\n"
    echo "scrapy crawl TweetScraper -a query=$line"
    scrapy crawl TweetScraper -a query="$line"
done < $file
