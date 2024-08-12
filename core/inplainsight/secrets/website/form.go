package website

import (
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"log"
)

func (s *WebsiteCredential) GetForm() *tview.Form {
	form := tview.NewForm()

	form.
		AddInputField("URL", s.GetTitle(), 0, nil, nil).
		AddInputField("Note", s.GetDescription(), 0, nil, nil).
		AddInputField("User", s.account, 0, nil, nil).
		AddPasswordField("Password", s.password, 0, '*', nil).
		AddButton("Cancel", func() {
			pages.GoBack()
		}).
		AddButton("Save", func() {
			s.website = form.GetFormItemByLabel("URL").(*tview.InputField).GetText()
			s.note = form.GetFormItemByLabel("Note").(*tview.InputField).GetText()
			s.account = form.GetFormItemByLabel("User").(*tview.InputField).GetText()
			s.password = form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

			err := inplainsight.Conceal(s)

			if err == nil {
				form.GetFormItemByLabel("URL").(*tview.InputField).SetText("")
				form.GetFormItemByLabel("Note").(*tview.InputField).SetText("")
				form.GetFormItemByLabel("User").(*tview.InputField).SetText("")
				form.GetFormItemByLabel("Password").(*tview.InputField).SetText("")
				inplainsight.InPlainSight.App.SetFocus(form.GetFormItem(0))

				log.Println("added secret", s)
			}

			pages.GoBack()
		}).
		SetButtonsAlign(tview.AlignRight)

	return form
}
