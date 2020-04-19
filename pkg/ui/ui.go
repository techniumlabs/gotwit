package ui

import (
	"github.com/rivo/tview"
	gotwitApp "github.com/techniumlabs/gotwit/pkg/app"
)

type AppUI struct {
	UIApp *tview.Application
}

func NewUI(app *gotwitApp.App) *AppUI {

	uiapp := tview.NewApplication()

	return &AppUI{UIApp: uiapp}

}

func (a *AppUI) Render() {
	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	if err := a.UIApp.SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}
