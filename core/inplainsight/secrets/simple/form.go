package simple

import (
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"github.com/zangarmarsh/inplainsight/core/utility/tuihelpers"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"log"
	"path/filepath"
)

func (s *SimpleSecret) GetForm() *tview.Form {
	form := tview.NewForm()

	form.
		AddInputField("Title", s.title, 0, nil, nil).
		AddInputField("Description", s.description, 0, nil, nil).
		AddPasswordField("Secret", s.secret, 0, '*', nil)

	var dedicatedHostInput *tview.InputField
	{
		var containerPath string

		if s.GetContainer() != nil {
			containerPath = s.GetContainer().Host.GetPath()
			containerPath = containerPath[len(inplainsight.InPlainSight.Path):]
		}

		dedicatedHostInput = tuihelpers.GenerateContainerSelector(containerPath)
		form.AddFormItem(dedicatedHostInput)
	}

	form.
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
