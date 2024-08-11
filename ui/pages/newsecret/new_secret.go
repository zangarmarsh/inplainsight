package newsecret

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets/simple"
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets/website"
	"github.com/zangarmarsh/inplainsight/ui/pages"
)

func GetName() string {
	return "new"
}

func Create() *pages.GridPage {
	var secretModel secrets.SecretInterface
	var listOfModels *tview.List

	page := pages.GridPage{}
	page.SetName(GetName())

	grid := tview.NewGrid().
		SetRows(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0).
		SetColumns(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

	flex := tview.NewFlex()
	flex.
		SetTitle(fmt.Sprintf("new - inplainsight v%s ", inplainsight.Version)).
		SetBorder(true)

	flex.SetDirection(tview.FlexColumn)
	flex.SetBorderPadding(0, 0, 1, 1)

	var form *tview.Form
	setSecretForm := func() {
		if form != nil {
			flex.RemoveItem(form)
		}

		form = (secretModel).GetForm()
		inplainsight.InPlainSight.App.SetFocus(form)

		form.
			SetBorder(true).
			SetBorderPadding(3, 2, 2, 2)

		form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if form.GetFormItem(0).HasFocus() {
				if event.Key() == tcell.KeyBacktab {
					inplainsight.InPlainSight.App.SetFocus(listOfModels)
					return nil
				}
			}

			return event
		})

		// Removing placeholder
		flex.RemoveItem(nil)
		flex.AddItem(form, 0, 1, true)
	}

	listOfModels = tview.NewList()

	listOfModels.
		SetBorder(true).
		SetBorderPadding(2, 2, 2, 2).
		SetTitle("What kind of secret are you saving?")

	listOfModels.
		AddItem("Simple", "Generic Title/Description/Secret", 0, func() {
			secretModel = &simple.SimpleSecret{}
			setSecretForm()
		})

	listOfModels.
		AddItem("Website", "URL/Note/Account/Password", 0, func() {
			secretModel = &website.WebsiteCredential{}
			setSecretForm()
		})

	listOfModels.
		AddItem("Exit", "Press to get back to the secrets list", 'q', func() {
			pages.GoBack()
		})

	flex.AddItem(listOfModels, 0, 1, true)

	// Adding form's placeholder
	flex.AddItem(nil, 0, 1, false)

	grid.
		AddItem(flex, 2, 2, 4, 8, 30, 200, true).
		AddItem(flex, 1, 1, 4, 10, 0, 0, true)

	page.SetPrimitive(grid)

	return &page
}
