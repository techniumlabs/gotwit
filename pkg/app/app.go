package app

import (
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

func (a *App) HomeTimeline(c chan *twitter.Tweet, command string) {
	go func() {

		var tweets []twitter.Tweet
		var err error

		switch command {
		case "init":
			log.Println("Init")
			// tweets, _, err = a.Client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
			//	Count: 20,
			// })
		case "scroll":
			log.Printf("scrolling %d tweets")

		}

		tweets, _, err = a.Client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
			Count: 100,
		})
		if err != nil {
			log.Error(err.Error())
		}
		for _, tweet := range tweets {
			c <- &tweet
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
