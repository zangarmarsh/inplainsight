package simple

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"log"
	"path/filepath"
)

func (s *SimpleSecret) GetForm() *tview.Form {
	form := tview.NewForm()

	dedicatedHostInput := tview.NewInputField()
	dedicatedHostInput.
		SetLabel("Container").
		SetAutocompleteFunc(func(currentText string) (entries []string) {
			if results := inplainsight.InPlainSight.Hosts.SearchByContainerPath(currentText); results != nil {
				for _, result := range results {
					// Todo check if secret can be contained by chosen container
					entries = append(entries, fmt.Sprintf("%s (cap %dM)", filepath.Base(result.Host.GetPath()), result.Host.Cap()/1e6))
				}
			}

			if len(entries) > 1 {
				entries = append([]string{"Random"}, entries...)
			}

			return
		})

	if s.GetContainer() != nil {
		dedicatedHostInput.SetText(s.GetContainer().Host.GetPath())
	}

	form.
		AddInputField("Title", s.title, 0, nil, nil).
		AddInputField("Description", s.description, 0, nil, nil).
		AddPasswordField("Secret", s.secret, 0, '*', nil).
		AddFormItem(dedicatedHostInput).
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

			if container := dedicatedHostInput.GetText(); container != "Random" {
				log.Printf("setting %s as container\n", filepath.Join(inplainsight.InPlainSight.Path, container))
				if container := inplainsight.InPlainSight.Hosts.SearchByContainerPath(filepath.Join(inplainsight.InPlainSight.Path, container)); len(container) == 1 {
					container := container[0]
					secrets.LinkSecretAndContainer(s, container)
				}
			}

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
