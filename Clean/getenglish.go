package main

import "github.com/abadojack/whatlanggo"

func detectEnglish(text string) bool {
	info := whatlanggo.Detect(text)

	if info.Lang.String() == "English" {
		// confidence has to be 51%

		if info.Confidence > 0.5 {
			return true
		}
	}

	return false
}
