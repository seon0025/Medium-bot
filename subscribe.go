package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

const (
    baseUrl = "https://medium.com/feed/tag"
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

func subscribe(category string, articleStream chan<- article, ticker *time.Ticker) {
    url := fmt.Sprintf("%s/%s", baseUrl, category)

    for tick := range ticker.C {
        func() {
            resp, err := http.Get(url)
            if err != nil {
                // might have to handle error
                panic(err)
            }

            defer resp.Body.Close()

            articles, err := parseArticles(resp.Body)
            if err != nil {
                panic(err)
            }

            took := time.Since(tick).Round(time.Millisecond)
            log.Printf("task with category: %s took: %d ms\n", category, int(took.Milliseconds()))

            for _, article := range articles {
                if _, exists := latestArticles[article.title]; !exists {
                    latestArticles[article.title] = time.Now().Add(time.Hour)
                    articleStream <- article
                }
            }
        }()
    }
}