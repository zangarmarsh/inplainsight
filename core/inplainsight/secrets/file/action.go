package file

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"github.com/zangarmarsh/inplainsight/ui/widgets"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (s *File) DoAction() {
	var commandList *tview.List

	page := &pages.GridPage{}
	page.SetName("show secret")

	grid := tview.NewGrid().
		SetRows(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0).
		SetColumns(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

	flex := tview.NewFlex()
	flex.
		SetBorder(true)

	flex.SetDirection(tview.FlexColumn)
	flex.SetBorderPadding(0, 0, 1, 1)

	grid.
		AddItem(flex, 2, 2, 8, 8, 30, 200, false).
		AddItem(flex, 1, 1, 8, 10, 0, 0, false)

	// Export to file form
	var exportForm *tview.Form
	{
		exportForm = tview.NewForm()

		pathInputField := tview.NewInputField()

		pathInputField.
			SetLabel("Path").
			SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyBacktab {
					inplainsight.InPlainSight.App.SetFocus(commandList)
					return nil
				}

				return event
			})

		automaticallyOpenCheckbox := tview.NewCheckbox()
		automaticallyOpenCheckbox.SetLabel("Open it up")

		exportForm.
			AddFormItem(pathInputField).
			AddFormItem(automaticallyOpenCheckbox).
			AddButton("Export", func() {
				outputFileName := pathInputField.GetText()

				// if it is a folder add the file name
				if filepath.Dir(outputFileName) == strings.TrimRight(outputFileName, string(os.PathSeparator)) {
					outputFileName = filepath.Join(outputFileName, s.fileName)
				}

				handle, err := os.Create(pathInputField.GetText())
				if err != nil {
					log.Println("Error creating file:", err)
					widgets.NewModal(widgets.ModalAlert, fmt.Sprintf("Cannot create file %s", outputFileName), "", nil)

					return
				}

				defer handle.Close()
				charCount, err := handle.WriteString(s.content)
				if err != nil || charCount != len(s.content) {
					log.Println("error writing to file:", err)
					widgets.NewModal(
						widgets.ModalAlert,
						fmt.Sprintf(
							"There was an unexpected error while writing content into %s",
							outputFileName,
						),
						"",
						nil,
					)

					return
				}

				if automaticallyOpenCheckbox.IsChecked() {
					err := exec.Command("open", outputFileName).Run()
					if err != nil {
						log.Println("Error opening file through `open`:", err)
						widgets.NewModal(
							widgets.ModalAlert,
							fmt.Sprintf(
								"There was an unexpected error while opening file %s",
								outputFileName,
							),
							"",
							nil,
						)

						return
					}
				}

			}).
			SetBorder(true).
			SetBorderPadding(3, 2, 2, 2)
	}

	{
		commandList = tview.NewList()
		commandList.
			SetBorder(true).
			SetBorderPadding(2, 2, 2, 2)

		commandList.
			AddItem("Export to file to disk", "", 'e', func() {
				// Do nothing, since it's the only tab it is open by default
				inplainsight.InPlainSight.App.SetFocus(exportForm)
			}).
			AddItem("Go back to secrets list", "", 'q', func() {
				pages.GoBack()
			})

		commandList.SetCurrentItem(0)
	}

	flex.AddItem(commandList, 0, 1, false)
	flex.AddItem(exportForm, 0, 1, true)

	page.SetPrimitive(grid)
	pages.Navigate(page)

	inplainsight.InPlainSight.App.SetFocus(commandList)
}
