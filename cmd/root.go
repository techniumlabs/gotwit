package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	gotwitApp "github.com/techniumlabs/gotwit/pkg/app"
	gotwitUI "github.com/techniumlabs/gotwit/pkg/ui"
)

var cfgFile, consumerKey, consumerToken, accessToken, accessSecret string
var app *gotwitApp.App

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotwit",
	Short: "Terminal Client for Twitter",
	Long:  `A Terminal Client and cli for twitter`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		app, err = gotwitApp.NewApp(cfgFile)
		if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Hello Twits")
		ui := gotwitUI.NewUI(app)
		ui.Render()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute: %s", err.Error())
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	viper.SetEnvPrefix("TWITTER")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotwit.yaml)")
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
