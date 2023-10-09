package main

import (
	"os"
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
	token := os.Getenv("DISCORD_TOKEN")
	channelId := os.Getenv("CHAN_ID")

	articleStream := make(chan article)
	go watchArticleExpired()

	client, err := createClient(token, channelId)
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
