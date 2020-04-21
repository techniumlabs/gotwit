package ui

import (
	"fmt"
	"io"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/gdamore/tcell"
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

	if err := a.UIApp.SetRoot(a.GetHomeLayout(), true).Run(); err != nil {
		panic(err)
	}
}

func (a *AppUI) GetHomeLayout() *tview.Flex {
	frame := tview.NewFrame(a.GetTimelineView()).
		SetBorders(2, 2, 2, 2, 4, 4).
		AddText("Home", true, tview.AlignLeft, tcell.ColorRed)

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(frame, 0, 1, true), 0, 1, false)

	return flex
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
		go a.App.HomeTimeline(c, "init")
		for t := range c {
			DisplayTweet(textView, t)
		}
	}()

	return textView
}

func DisplayTweet(w io.Writer, t *twitter.Tweet) {
	fmt.Fprintf(w, "[green]%s[white] @%s [red]❤️[white] %d  RT%d\n", t.User.Name, t.User.ScreenName, t.FavoriteCount, t.RetweetCount)
	fmt.Fprintf(w, "%s\n\n", t.Text)
}
