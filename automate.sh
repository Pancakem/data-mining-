#!/bin/bash

file="Parameters/users.txt"
echo "Starting crawling.."
while read line; do
    echo $line
    python3 instagram-crawler/crawler.py posts_full -u $line -o "Results/instagram_data/$line.json"
done <  $file