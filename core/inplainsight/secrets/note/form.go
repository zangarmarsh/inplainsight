package note

import (
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"github.com/zangarmarsh/inplainsight/core/utility/tuihelpers"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"log"
	"path/filepath"
)

func (s *Note) GetForm() *tview.Form {
	form := tview.NewForm()

	form.
		AddInputField("Title", s.title, 0, nil, nil).
		AddTextArea("Note", s.note, 0, 0, 0, nil).
		AddCheckbox("Hidden by default?", s.isHidden, nil)

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
			s.title = form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
			s.isHidden = form.GetFormItemByLabel("Hidden by default?").(*tview.Checkbox).IsChecked()
			s.note = form.GetFormItemByLabel("Note").(*tview.TextArea).GetText()

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
				form.GetFormItemByLabel("Hidden by default?").(*tview.Checkbox).SetChecked(false)
				form.GetFormItemByLabel("Note").(*tview.TextArea).SetText("", false)

				inplainsight.InPlainSight.App.SetFocus(form.GetFormItem(0))

				log.Println("added secret", s)
			}

			pages.GoBack()
		}).
		SetButtonsAlign(tview.AlignRight)

	return form
}
