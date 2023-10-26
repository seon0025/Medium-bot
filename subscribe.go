package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	baseUrl = "https://blog.msg-team.com"
)

var (
	latestArticles = make(map[string]time.Time)
)

func watchArticleExpired() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		now := time.Now()

		for title, expiresAt := range latestArticles {
			if now.After(expiresAt) {
				delete(latestArticles, title)
			}
		}
	}
}

func doSubscribe(url string, articleStream chan<- article) {
	resp, err := http.Get(url)
	if err != nil {
		// might have to handle error
		panic(err)
	}

	defer resp.Body.Close()

	articles, err := parseArticles(resp.Body)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	for _, article := range articles {
		if _, exists := latestArticles[article.title]; !exists {
			latestArticles[article.title] = time.Now().Add(time.Hour)
			articleStream <- article
		}
	}
}

func subscribe(category string, articleStream chan<- article, ticker *time.Ticker) {
	url := fmt.Sprintf("%s/%s", baseUrl, category)

	doSubscribe(url, articleStream)

	for tick := range ticker.C {
		doSubscribe(url, articleStream)
		took := time.Since(tick).Round(time.Millisecond)
		log.Printf("task with category: %s took: %d ms\n", category, int(took.Milliseconds()))
	}
}
