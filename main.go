package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"github.com/kelseyhightower/envconfig"
	irc "github.com/thoj/go-ircevent"
)

const serverssl = "irc.chat.twitch.tv:6697"

type TwitchConfig struct {
	Nick     string
	Password string
}

type TwitterConfig struct {
	ConsumerKey       string `envconfig:"CONSUMER_KEY"`
	ConsumerSecret    string `envconfig:"CONSUMER_SECRET"`
	AccessToken       string `envconfig:"ACCESS_TOKEN"`
	AccessTokenSecret string `envconfig:"ACCESS_TOKEN_SECRET"`
}

func main() {
	ch1 := make(chan bool)
	ch2 := make(chan bool)

	go startTwitchCommentStream(ch1)
	go startTwitterHashTagStream(ch2)

	<-ch1
	<-ch2
}

func startTwitchCommentStream(done chan bool) {
	var config TwitchConfig
	envconfig.Process("TWITCH", &config)

	nick := config.Nick
	con := irc.IRC(nick, nick)

	con.Password = config.Password
	con.UseTLS = true
	con.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	channel := os.Args[1]
	con.AddCallback("001", func(e *irc.Event) { con.Join(channel) })
	con.AddCallback("PRIVMSG", printTwitchMessage)
	err := con.Connect(serverssl)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	go con.Loop()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("quit twitch irc...")
	con.Quit()
	done <- true
}

func startTwitterHashTagStream(done chan bool) {
	var c TwitterConfig
	envconfig.Process("TWITTER", &c)
	config := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	demux := twitter.NewSwitchDemux()
	demux.Tweet = printTweet

	track := os.Args[2]
	fmt.Printf("start twitter streaming: %s\n", track)
	filterParams := &twitter.StreamFilterParams{Track: []string{track}}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	quitTwitter := make(chan os.Signal)
	signal.Notify(quitTwitter, syscall.SIGINT, syscall.SIGTERM)
	<-quitTwitter

	fmt.Println("quit twitter stream...")
	stream.Stop()
	done <- true
}

func printTweet(tweet *twitter.Tweet) {
	fmt.Printf("<twitter> %s: %s\n", tweet.User.ScreenName, tweet.Text)
}

func printTwitchMessage(e *irc.Event) {
	fmt.Printf("[twitch] %s: %s\n", e.User, e.Arguments[1])
}
