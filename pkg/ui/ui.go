package ui

import (
	"bytes"
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
	home := a.GetHomeLayout()
	a.UIApp.SetRoot(home, true)
	a.UIApp.SetFocus(home)
	if err := a.UIApp.Run(); err != nil {
		panic(err)
	}
}

func (a *AppUI) GetHomeLayout() *tview.Flex {
	table := a.GetTimelineViewAsTable()
	frame := tview.NewFrame(table).
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

func (a *AppUI) GetTimelineViewAsTable() *tview.Table {
	log.Println("Setting Timeline View")
	table := tview.NewTable()
	table.SetBorder(false)

	go func() {
		c := make(chan *twitter.Tweet)
		go a.App.HomeTimeline(c, "init")
		index := 0
		for t := range c {
			table.SetCell(index, 0, tview.NewTableCell(FormatTweet(t)))
			index += 1
		}
	}()

	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			a.UIApp.Stop()
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
		log.Println("key pressed")
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(false, false)
	})
	return table
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
