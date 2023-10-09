package main

import (
    "flag"
    "time"
)

const (
    duration = 10 * time.Second
)

var categories = []string{"programming", "coding"}

type article struct {
    title string
    link  string
}

func main() {
    token := flag.String("token", "", "your discord bot's token")
    channelId := flag.String("chanId", "", "your notification channel id")

    flag.Parse()

    articleStream := make(chan article)
    go watchArticleExpired()

    client, err := createClient(*token, *channelId)
    if err != nil {
        panic(err)
    }
    defer client.Close()

    ticker := time.NewTicker(duration)

    for _, category := range categories {
        go subscribe(category, articleStream, ticker)
    }

    for article := range articleStream {
        go notify(client, article)
    }
}