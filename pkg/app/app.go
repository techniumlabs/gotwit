package app

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/techniumlabs/gotwit/pkg/config"
)

type App struct {
	Config *config.Config
	client *twitter.Client
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

	return &App{Config: c, client: client}, nil
}
