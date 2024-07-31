package ui

import (
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"golang.design/x/clipboard"
)

func Bootstrap() {
	err := clipboard.Init()
	if err != nil {
		return
	}

	inplainsight.InPlainSight.App = tview.NewApplication()
	inplainsight.InPlainSight.Pages = tview.NewPages()
	inplainsight.InPlainSight.App.SetRoot(inplainsight.InPlainSight.Pages, true).EnableMouse(true)
}
