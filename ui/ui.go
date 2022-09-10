package ui

import (
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui/events"
	"golang.design/x/clipboard"
	"log"
	"strings"
)

type InPlainSightClient struct {
	events.EventeableStruct

	App            *tview.Application
	Pages          *tview.Pages

	InvolvedFiles  []string
	Secrets        []*Secret

	MasterPassword string
	Path           string
}

var InPlainSight = &InPlainSightClient{}

const (
	Version = "1.0.0"
)

func Bootstrap() {
	err := clipboard.Init()
	if err != nil {
		return
	}

	InPlainSight.App   = tview.NewApplication()
	InPlainSight.Pages = tview.NewPages()

	InPlainSight.App.SetRoot(InPlainSight.Pages, true).EnableMouse(true)
}
