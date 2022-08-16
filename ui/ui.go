package ui

import (
	"github.com/rivo/tview"
)

type InPlainSightClient struct {
	App   *tview.Application
	Pages *tview.Pages
}

var InPlainSight = &InPlainSightClient{}

const (
	Version = "1.0.0"
)

func Bootstrap() {
	InPlainSight.App   = tview.NewApplication()
	InPlainSight.Pages = tview.NewPages()

	InPlainSight.App.SetRoot(InPlainSight.Pages, true)
}
