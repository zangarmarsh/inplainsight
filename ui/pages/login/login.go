package login

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"github.com/zangarmarsh/inplainsight/ui/widgets"
)

type Page struct {
	pages.GridPage
}

type pageFactory struct { }

func (r pageFactory) GetName() string {
	return "login"
}

func (r pageFactory) Create() pages.PageInterface {
	page := Page{}
	page.SetName(r.GetName())


	form := tview.NewForm()

	form.
		SetTitle(fmt.Sprintf(" login - inplainsight v%s ", ui.Version)).
		SetBorder(true)

	form.
		AddPasswordField("Password", "", 0, '*', nil).
		SetButtonsAlign(tview.AlignCenter).
		AddButton("Enter", func() {
			text := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			modal := widgets.ModalError(text)

			if len(text) != 0 {
				ui.InPlainSight.Pages.AddPage("modal-error", modal, true, true)
			}
		}).
		AddButton("Quit (CTRL + C)", func() {
			ui.InPlainSight.App.Stop()
		})

	grid := tview.NewGrid().
		SetRows(0, 0, 0).
		SetColumns(0, 0, 0)

	grid.
		AddItem(form, 0, 0, 3, 3, 0, 0, true).
		AddItem(form, 1, 1, 1, 1, 32, 70, true)

	page.SetPrimitive(grid)

	return &page
}

func register(records pages.PageFactoryDictionary) pages.PageFactoryInterface {
	records["login"] = pageFactory{}

	return records["login"]
}
var _ = register(pages.PageFactories)