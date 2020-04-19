package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
)

type AuthConfig struct {
	ConsumerKey   string `mapstructure: "consumerkey"`
	ConsumerSecret string `mapstructure: "consumersecret"`
	AccessToken   string `mapstructure: "accesstoken"`
	AccessSecret  string `mapstructure: "accesssecret"`
}

type Config struct {
	AuthConfig AuthConfig `mapstructure:"auth"`
}

func Load(cfgFile string) (*Config, error) {
	var err error
	v := viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		v.AddConfigPath(".")
		v.AddConfigPath(home + "/")
		v.SetConfigName(".gotwit")
	}

	v.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err = v.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", v.ConfigFileUsed())
		c := &Config{}
		err = v.Unmarshal(c)
		return c, err
	} else {
		log.Warnf("%s", err.Error())
		return nil, err
	}
}
