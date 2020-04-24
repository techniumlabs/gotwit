package app

import (
	"math"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/techniumlabs/gotwit/pkg/config"
)

type App struct {
	Config *config.Config
	Client *twitter.Client
}

func NewApp(cfgFile string) (*App, error) {
	c, err := config.Load(cfgFile)
	if err != nil {
		return nil, err
	}

	oconfig := oauth1.NewConfig(c.AuthConfig.ConsumerKey, c.AuthConfig.ConsumerSecret)
	token := oauth1.NewToken(c.AuthConfig.AccessToken, c.AuthConfig.AccessSecret)
	httpClient := oconfig.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return &App{Config: c, Client: client}, nil
}

func (a *App) HomeTimeline(output chan *twitter.Tweet, command chan string) {
	go func() {

		var tweets []twitter.Tweet
		var err error
		var maxTweet int64
		var minTweet int64 = math.MaxInt64

		for {
			switch <-command {
			case "init":
				tweets, _, err = a.Client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
					Count: 100,
				})
				if err != nil {
					log.Error(err.Error())
				}
			case "refresh":
				tweets, _, err = a.Client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
					Count: 100,
				})
				if err != nil {
					log.Error(err.Error())
				}
			case "next":
				tweets, _, err = a.Client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
					MaxID: minTweet,
					Count: 100,
				})
				if err != nil {
					log.Error(err.Error())
				}
			}

			for _, tweet := range tweets {
				if maxTweet < tweet.ID {
					maxTweet = tweet.ID
				}

				if minTweet > tweet.ID {
					minTweet = tweet.ID
				}
				output <- &tweet
			}
		}

	}()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

}

func (a *App) UserTweets(c chan *twitter.Tweet) {

	params := &twitter.StreamUserParams{
		StallWarnings: twitter.Bool(true),
	}
	stream, err := a.Client.Streams.User(params)
	// params := &twitter.StreamSampleParams{
	//	StallWarnings: twitter.Bool(true),
	// }
	// stream, err := a.Client.Streams.Sample(params)
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		c <- tweet
	}
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	stream.Stop()
}
