package simple

import (
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"log"
)

func (s *SimpleSecret) GetForm() *tview.Form {
	form := tview.NewForm()

	form.
		AddInputField("Title", s.title, 0, nil, nil).
		AddInputField("Description", s.description, 0, nil, nil).
		AddPasswordField("Secret", s.secret, 0, '*', nil).
		AddButton("Cancel", func() {
			pages.GoBack()
		}).
		AddButton("Save", func() {
			formTitle := form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
			formDescription := form.GetFormItemByLabel("Description").(*tview.InputField).GetText()
			formSecret := form.GetFormItemByLabel("Secret").(*tview.InputField).GetText()

			s.title = formTitle
			s.description = formDescription
			s.secret = formSecret

			err := inplainsight.Conceal(s)

			if err == nil {
				form.GetFormItemByLabel("Title").(*tview.InputField).SetText("")
				form.GetFormItemByLabel("Description").(*tview.InputField).SetText("")
				form.GetFormItemByLabel("Secret").(*tview.InputField).SetText("")
				inplainsight.InPlainSight.App.SetFocus(form.GetFormItem(0))

				log.Println("added secret", s)

				pages.GoBack()
			}
		}).
		SetButtonsAlign(tview.AlignRight)

	return form
}
