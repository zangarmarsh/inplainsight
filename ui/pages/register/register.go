package register

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"log"
)

type Page struct {
	pages.GridPage
}

type pageFactory struct { }

func (r pageFactory) GetName() string {
	return "register"
}

func (r pageFactory) Create() pages.PageInterface {
	page := Page{}
	page.SetName(r.GetName())

	form := tview.NewForm()

	form.
		SetBorder(false)

	form.
		AddPasswordField("Master Password", "", 0, '*', nil).
		SetButtonsAlign(tview.AlignCenter).
		AddButton("Register", func() {
			text := form.GetFormItemByLabel("Master Password").(*tview.InputField).GetText()
			log.Print(text)

		}).
		AddButton("Quit (CTRL + C)", func() {
			ui.InPlainSight.App.Stop()
		})

	grid := tview.NewGrid().
		SetRows(0, 0).
		SetColumns(0, 0, 0)

	flex := tview.NewFlex()
	flex.SetTitle(fmt.Sprintf(" register - inplainsight v%s ", ui.Version)).
			 SetBorder(true)

	flex.SetDirection(tview.FlexRow)
	flex.SetBorderPadding(2,2,2,2)

	text := tview.NewTextView().
					      SetText("Greetings, stranger. To get it started you'll just have to set a strong password.\n\n[red::b]Beware:[-:-:-] don't forget it or you won't be able to access to your secrets never again.").
								SetWordWrap(true).
								SetDynamicColors(true)

	flex.AddItem(text, 0, 2, false)
	flex.AddItem(form, 0, 1, true)

	grid.
		AddItem(flex, 0, 1, 1, 1, 30, 50, true).
		AddItem(flex, 0,0 ,3, 3, 0, 0, true)

	grid.SetBorderPadding(2, 0, 0, 0)

	page.SetPrimitive( grid )

	return &page
}

func register(records pages.PageFactoryDictionary) pages.PageFactoryInterface {
	records["register"] = pageFactory{}
	return records["register"]
}
var _ = register(pages.PageFactories)
