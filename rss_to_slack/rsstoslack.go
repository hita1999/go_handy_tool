package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

// RSS structure
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title string `xml:"title"`
		Link  string `xml:"link"`
		Item  []struct {
			Title       string `xml:"title"`
			Description string `xml:"description"`
			Link        string `xml:"link"`
		} `xml:"item"`
	} `xml:"channel"`
}


func main() {
	// RSS URL
	rssURL := "https://zenn.dev/topics/python/feed"

	// Get RSS
	resp, err := http.Get(rssURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Transform response to []byte
	var body strings.Builder
	if _, err := io.Copy(&body, resp.Body); err != nil {
		panic(err)
	}

	// Parse RSS
	var rss RSS
	if err := xml.Unmarshal([]byte(body.String()), &rss); err != nil {
		panic(err)
	}

	// Slack Channel
	channel := "#news_test"

	// Slack API Token
	apiToken := os.Getenv("SLACK_API_TOKEN")

	// Initialize Slack Client
	client := slack.New(apiToken)

	// Make Message
	var messages []string
	for _, item := range rss.Channel.Item {
		messages = append(messages, fmt.Sprintf("<%s|%s>\n%s", item.Link, item.Title, item.Description))
	}
	message := strings.Join(messages, "\n")

	// Post Message
	_, _, err = client.PostMessage(channel, slack.MsgOptionText(message, false))
	if err != nil {
		panic(err)
	}

	fmt.Println("Post to Slack channel: ", channel)
}