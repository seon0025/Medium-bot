package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var channelId string

func createClient(token, chanId string) (*discordgo.Session, error) {
	client, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}
	log.Println(token, chanId)

	client.Identify.Intents = discordgo.IntentGuildMessages

	err = client.Open()
	if err != nil {
		return nil, err
	}

	channelId = chanId

	return client, nil
}

func notify(session *discordgo.Session, article article) {

	//content := fmt.Sprintf("%#v, here", article)
	content := fmt.Sprintf("📰 **%s**\n%s", article.title, article.link)

	_, err := session.ChannelMessageSend(channelId, content, func(cfg *discordgo.RequestConfig) {
		cfg.Request.Header.Set("Authorization", "Bot "+session.Identify.Token)
	})
	if err != nil {
		panic(err)
	}
}
