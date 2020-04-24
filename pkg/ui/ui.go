package ui

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

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
	home := a.GetHomeLayout()
	a.UIApp.SetRoot(home, true)
	a.UIApp.SetFocus(home)
	if err := a.UIApp.Run(); err != nil {
		panic(err)
	}
}

func (a *AppUI) GetHomeLayout() *tview.Frame {
	table := a.GetTimelineView()
	frame := tview.NewFrame(table).
		SetBorders(2, 2, 2, 2, 4, 4).
		AddText("Home", true, tview.AlignLeft, tcell.ColorRed)

	// flex := tview.NewFlex().
	//	AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
	//		AddItem(frame, 0, 1, true), 0, 1, false)

	return frame
}

func (a *AppUI) GetTimelineView() *tview.TextView {
	log.Println("Setting Timeline View")
	totalTweets := 0
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			a.UIApp.Draw()
		})

	c := make(chan *twitter.Tweet)
	cin := make(chan string)
	go func() {
		go a.App.HomeTimeline(c, cin)
		cin <- "init"
		for t := range c {
			fmt.Fprintf(textView, `["%d"]`, totalTweets)
			DisplayTweet(textView, t)
			fmt.Fprintf(textView, `[""]`)
			totalTweets += 1
		}
	}()

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		switch key {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				currentSelection := textView.GetHighlights()
				if len(currentSelection) > 0 {
					next, _ := strconv.Atoi(currentSelection[0])
					if next < totalTweets-1 {
						next += 1
					} else {
						cin <- "next"
					}
					textView.Highlight(strconv.Itoa(next)).ScrollToHighlight()
				} else {
					textView.Highlight("0").ScrollToHighlight()
				}
			case 'k':
				currentSelection := textView.GetHighlights()
				if len(currentSelection) > 0 {
					next, _ := strconv.Atoi(currentSelection[0])
					if next > 0 {
						next -= 1
					} else {
						cin <- "refresh"
					}
					textView.Highlight(strconv.Itoa(next)).ScrollToHighlight()
				} else {
					textView.Highlight("0").ScrollToHighlight()
				}
			}
		}
		return event
	})
	return textView
}

func DisplayTweet(w io.Writer, t *twitter.Tweet) {
	fmt.Fprintf(w, "[green]%s[white] @%s [red]❤️[white] %d  RT%d\n", t.User.Name, t.User.ScreenName, t.FavoriteCount, t.RetweetCount)
	fmt.Fprintf(w, "%s\n\n", t.Text)
}

func FormatTweet(t *twitter.Tweet) string {
	var result bytes.Buffer

	result.WriteString(fmt.Sprintf("[green]%s[white] @%s [red]❤️[white] %d  RT%d\n", t.User.Name, t.User.ScreenName, t.FavoriteCount, t.RetweetCount))
	result.WriteString(fmt.Sprintf("%s\n\n", t.Text))
	return result.String()
}
