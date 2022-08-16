package list

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui"
	"github.com/zangarmarsh/inplainsight/ui/pages"
)

type Page struct {
	pages.GridPage
}

type pageFactory struct { }

func (r pageFactory) GetName() string {
	return "list"
}

func (r pageFactory) Create() pages.PageInterface {
	page := Page{}
	r.GetName()

	grid := tview.NewGrid().
		SetRows(0, 0).
		SetColumns(0, 0, 0)

	flex := tview.NewFlex()
	flex.SetTitle(fmt.Sprintf(" inplainsight v%s ", ui.Version)).
		SetBorder(true)

	// @ToDo: Implement here the interface

	flex.SetDirection(tview.FlexRow)
	page.SetPrimitive(grid)

	return &page
}


func register(records pages.PageFactoryDictionary) pages.PageFactoryInterface {
	records["list"] = pageFactory{}

	return records["list"]
}
var _ = register(pages.PageFactories)