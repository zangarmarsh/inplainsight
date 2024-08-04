package newsecret

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"log"
)

func GetName() string {
	return "new"
}

func Create() *pages.GridPage {
	page := pages.GridPage{}
	page.SetName(GetName())

	form := tview.NewForm()

	form.
		SetBorder(false)

	form.
		AddInputField("Title", "", 0, nil, nil).
		AddInputField("Description", "", 0, nil, nil).
		AddPasswordField("Container", "", 0, '*', nil).
		SetButtonsAlign(tview.AlignCenter).
		AddButton("Save", func() {
			formTitle := form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
			formDescription := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
			formSecret := form.GetFormItemByLabel("Container").(*tview.InputField).GetText()

			err := pages.Navigate("list")
			if err != nil {
				// Todo handle it
				log.Fatal(err)
			}

			secret := inplainsight.Secret{}

			secret.Title = formTitle
			secret.Description = formDescription
			secret.Secret = formSecret

			err = inplainsight.Conceal(&secret)

			if err == nil {
				form.GetFormItemByLabel("Title").(*tview.InputField).SetText("")
				form.GetFormItemByLabel("Description").(*tview.InputField).SetText("")
				form.GetFormItemByLabel("Container").(*tview.InputField).SetText("")
				inplainsight.InPlainSight.App.SetFocus(form.GetFormItem(0))

				log.Println("added secret", secret)
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
