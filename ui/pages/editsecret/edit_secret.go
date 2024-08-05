package editsecret

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"log"
)

func GetName() string {
	return "edit"
}

func Create(secret *inplainsight.Secret) *pages.GridPage {
	page := pages.GridPage{}
	page.SetName(GetName())

	form := tview.NewForm()

	form.
		SetBorder(false)

	var formTitle string
	var formDescription string
	var formSecret string

	if secret != nil {
		formTitle = secret.Title
		formDescription = secret.Description
		formSecret = secret.Secret
	}

	form.
		AddInputField("Title", formTitle, 0, nil, nil).
		AddInputField("Description", formDescription, 0, nil, nil).
		AddPasswordField("Secret", formSecret, 0, '*', nil).
		SetButtonsAlign(tview.AlignCenter).
		AddButton("Save", func() {
			formTitle = form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
			formDescription = form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
			formSecret = form.GetFormItemByLabel("Secret").(*tview.InputField).GetText()

			err := pages.Navigate("list")

			(*secret).Title = formTitle
			(*secret).Description = formDescription
			(*secret).Secret = formSecret

			err = inplainsight.Conceal(secret)

			if err == nil {
				log.Println("secret updated", secret)
			} else {
				log.Println("error while concealing", err)
			}

			inplainsight.InPlainSight.Pages.RemovePage(GetName())
		}).
		AddButton("Back", func() {
			inplainsight.InPlainSight.Pages.RemovePage(GetName())
		})

	grid := tview.NewGrid().
		SetRows(0, 0).
		SetColumns(0, 0, 0)

	flex := tview.NewFlex()
	flex.SetTitle(fmt.Sprintf(" new - inplainsight v%s ", inplainsight.Version)).
		SetBorder(true)

	flex.SetDirection(tview.FlexRow)
	flex.SetBorderPadding(2, 2, 2, 2)

	flex.AddItem(form, 0, 1, true)

	grid.
		AddItem(flex, 0, 1, 1, 1, 30, 50, true).
		AddItem(flex, 0, 0, 3, 3, 0, 0, true)

	grid.SetBorderPadding(2, 0, 0, 0)

	page.SetPrimitive(grid)

	return &page
}
