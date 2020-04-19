package ui

import (
	"fmt"
	"io"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
	gotwitApp "github.com/techniumlabs/gotwit/pkg/app"
)

type AppUI struct {
	UIApp *tview.Application
	App   *gotwitApp.App
}

func NewUI(app *gotwitApp.App) *AppUI {

	uiapp := tview.NewApplication()

	return &AppUI{UIApp: uiapp,
		App: app}

}

func (a *AppUI) Render() {
	// box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")

	if err := a.UIApp.SetRoot(a.GetTimelineView(), true).Run(); err != nil {
		panic(err)
	}
}

func (a *AppUI) GetTimelineView() *tview.TextView {
	log.Println("Setting Timeline View")
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			a.UIApp.Draw()
		})

	go func() {
		c := make(chan *twitter.Tweet)
		go a.App.UserTweets(c)
		for t := range c {
			DisplayTweet(textView, t)
		}
	}()

	textView.SetBorder(true)
	return textView
}

func DisplayTweet(w io.Writer, t *twitter.Tweet) {
	fmt.Fprintf(w, "[red]%s[white]\n", t.User.Name)
	fmt.Fprintf(w, "%s\n\n", t.Text)
}
