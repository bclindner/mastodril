package main

import (
	"context"
	"github.com/mattn/go-mastodon"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/go-twitter/twitter"
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"
	"encoding/binary"
	"os"
)

type Config struct {
	Username string `json:"handle"`
	ExcludeReplies *bool `json:"exclude_replies"`
	Twitter TwitterConfig `json:"twitter"`
	Mastodon MastodonConfig `json:"mastodon"`
}

type TwitterConfig struct {
	ConsumerKey string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
	AccessToken string `json:"accessToken"`
	AccessSecret string `json:"accessSecret"`
}

type MastodonConfig struct {
	Server string `json:"server"`
	ClientID string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	AccessToken string `json:"accessToken"`
}

func GetTwitterClient (config TwitterConfig) *twitter.Client {
	// setup oauth1 stuff
	oauth := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret)
	token := oauth1.NewToken(config.AccessToken, config.AccessSecret)
	// create the client
	httpclient := oauth.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpclient)
	// return it
	return client
}

func GetMastodonClient(config MastodonConfig) *mastodon.Client {
	// map MastodonConfig to an actual proper mastodon.Config, run NewClient with it
	client := mastodon.NewClient(&mastodon.Config{
		Server: config.Server,
		ClientID: config.ClientID,
		ClientSecret: config.ClientSecret,
		AccessToken: config.AccessToken,
	})
	// return the result
	return client
}

func main () {
	// get the last checked tweet
	var lastTweetID int64 = 0
	lastTweetFile, err := ioutil.ReadFile("mastodril.last")
	if err == nil {
		lastTweetID, _ = binary.Varint(lastTweetFile)
	}
	// get config file
	configFile, err := ioutil.ReadFile("mastodril.json")
	if err != nil { fmt.Println(err); os.Exit(1) }
	// parse the config
	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil { fmt.Println(err); os.Exit(1) }
	// get api clients
	t := GetTwitterClient(config.Twitter)
	m := GetMastodonClient(config.Mastodon)
	// get twitter timeline starting at the user in question
	tweets, _, err := t.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: config.Username,
		SinceID: lastTweetID,
		Count: 5,
		ExcludeReplies: config.ExcludeReplies,
		TweetMode: "extended",
	})
	if err != nil { fmt.Println(err); os.Exit(1) }
	// post each tweet since the last checked tweet ID
	// done in reverse order to keep chronological order
	for i := len(tweets)-1; i >= 0; i-- {
		tweet := tweets[i]
		fmt.Println("found new tweet with id", tweet.IDStr)
		fulltext := tweet.FullText
		for _, e := range tweet.Entities.Urls {
			expandedurl := strings.ReplaceAll(e.ExpandedURL, "http://", "")
			fulltext = strings.ReplaceAll(fulltext, e.URL, expandedurl)
		}
		fulltext = strings.ReplaceAll(fulltext, "&amp;", "&")
		fulltext = strings.ReplaceAll(fulltext, "&lt;", "<")
		fulltext = strings.ReplaceAll(fulltext, "&gt;", ">")
		m.PostStatus(context.Background(), &mastodon.Toot{
			Status: fulltext,
		})
		lastTweetID = tweet.ID
	}
	// write the last tweet ID to the .last file
	lastIDBinary := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(lastIDBinary, lastTweetID)
	err = ioutil.WriteFile("mastodril.last", lastIDBinary, os.FileMode(int(0644)))
	if err != nil { fmt.Println(err); os.Exit(1) }
	// done!
}
