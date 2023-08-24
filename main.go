package main

import (
	"fmt"
	"net/http"
)

type requestResult struct {
	url string
	status string
}

func main() {
	results := make(map[string]string)
	c := make(chan requestResult)

	urls := []string{
		"https://www.google.com/",
		"https://www.aws.com/",
		"https://www.airbnb.com/",
		"https://www.reddit.com/",
		"https://www.naver.com/",
		"https://www.daum.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
	}

	for _, url := range urls {
		go hitURL(url, c)
	}
	for i := 0; i < len(urls); i++ {
		result := <-c
		results[result.url] = result.status
	}
	for url, status := range results {
		fmt.Println(url, status)
	}
}

func hitURL(url string, c chan<- requestResult) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		c <- requestResult{url:url, status: "FAILED"}
	} else {
		c <- requestResult{url:url, status: "OK"}
	}
}