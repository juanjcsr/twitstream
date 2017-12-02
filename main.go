package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	params := &twitter.StreamFilterParams{
		//Track:         []string{""},
		Locations:     []string{"-118.6523001,14.3895,-86.5887,32.7186534"},
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.Filter(params)
	if err != nil {
		log.Panicf("could not read stream: %v", err)
	}

	f, err := os.OpenFile("stream.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Panicf("could not open file: %v", err)
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(t *twitter.Tweet) {
		str, _ := json.Marshal(t)
		_, err := f.WriteString(string(str))
		f.WriteString("\r\n")
		if err != nil {
			log.Panicf("could not write to file: %v", err)
		}
	}

	demux.HandleChan(stream.Messages)
}
