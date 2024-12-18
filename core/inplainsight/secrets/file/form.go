package file

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"github.com/zangarmarsh/inplainsight/core/utility"
	"github.com/zangarmarsh/inplainsight/core/utility/tuihelpers"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"github.com/zangarmarsh/inplainsight/ui/widgets"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func (s *File) GetForm() *tview.Form {
	var filePathInput *tview.InputField
	form := tview.NewForm()

	filePathInput = tview.NewInputField()
	filePathInput.
		SetLabel("File path").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				if prevision := utility.SuggestFSPath(filePathInput.GetText()); prevision != "" {
					if prevision != filePathInput.GetText() {
						filePathInput.SetText(prevision)
						return nil
					}
				} else {
					return nil
				}
			}

			return event
		})

	form.
		AddInputField("Title", s.title, 0, nil, nil).
		AddTextArea("Note", s.note, 0, 0, 0, nil).
		AddFormItem(filePathInput)

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
			filePath := filePathInput.GetText()

			if filePath[0] == '~' && (runtime.GOOS == "linux" || runtime.GOOS == "darwin") {
				if homeDir, err := os.UserHomeDir(); err == nil {
					filePath = filepath.Join(homeDir, filePath[1:])
				}
			}

			var stat os.FileInfo
			var err error
			if stat, err = os.Stat(filePath); err != nil || stat == nil || !stat.Mode().IsRegular() {
				// Todo handle visually the error
				log.Printf("File %s is not a file", filePath)
				return
			}

			s.title = form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
			s.note = form.GetFormItemByLabel("Note").(*tview.TextArea).GetText()
			s.fileName = filepath.Base(filePath)

			if content, err := os.ReadFile(filePath); err != nil || int64(len(content)) != stat.Size() {
				// Todo handle the error
				log.Printf("Could not read (entirely?) file: %s (err: %s)", filePath, err)
				return
			} else {
				s.content = content
			}

			if container := dedicatedHostInput.GetText(); container != "Random" {
				log.Printf("setting %s as container\n", filepath.Join(inplainsight.InPlainSight.Path, container))
				if container := inplainsight.InPlainSight.Hosts.SearchByContainerPath(filepath.Join(inplainsight.InPlainSight.Path, container)); len(container) == 1 {
					container := container[0]
					secrets.LinkSecretAndContainer(s, container)
				}
			}

			err = inplainsight.Conceal(s)

			if err == nil {
				form.GetFormItemByLabel("Title").(*tview.InputField).SetText("")
				form.GetFormItemByLabel("Note").(*tview.TextArea).SetText("", false)
				filePathInput.SetText("")

				inplainsight.InPlainSight.App.SetFocus(form.GetFormItem(0))

				log.Println("added secret", s)
				pages.GoBack()
			} else {
				widgets.NewModal(widgets.ModalAlert, err.Error(), "", nil)
				log.Printf("Could not conceal file %s (err: %s)", filePath, err)
			}
		}).
		SetButtonsAlign(tview.AlignRight)

	return form
}
