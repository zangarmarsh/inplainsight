package editsecret

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"github.com/zangarmarsh/inplainsight/ui"
	"github.com/zangarmarsh/inplainsight/ui/events"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"github.com/zangarmarsh/inplainsight/ui/widgets"
	"log"
	"os"
	"strings"
	"time"
)

func GetName() string {
	return "edit"
}

func Create(secret *ui.Secret) *pages.GridPage {
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

			log.Println("reading folder", ui.InPlainSight.Path)
			files, err := os.ReadDir(ui.InPlainSight.Path)
			if err != nil {
				log.Println(err)
				widgets.ModalError(err.Error())
				ui.InPlainSight.App.ForceDraw()
				return
			}

			if len(files) == 0 {
				log.Println("empty directory")
				widgets.ModalError("The directory is empty")
				ui.InPlainSight.App.ForceDraw()
				return
			}

			err = pages.Navigate("list")

			s := steganography.Steganography{}

			for _, file := range files {
				if !file.IsDir() && strings.Contains(file.Name(), ".png") {
					filePath := fmt.Sprintf("%s/%s", strings.TrimRight(ui.InPlainSight.Path, "/\\"), file.Name())

					if filePath != secret.FilePath {
						continue
					}

					log.Println(
						"Concealing file ",
						fmt.Sprintf("%s/%s", strings.TrimRight(ui.InPlainSight.Path, "/\\"), file.Name()),
						fmt.Sprintf("with pass %#v", ui.InPlainSight.MasterPassword),
					)

					(*secret).Title = formTitle
					(*secret).Description = formDescription
					(*secret).Secret = formSecret

					err = s.Conceal(
						filePath,
						filePath,
						[]byte(secret.Serialize()),
						[]byte(ui.InPlainSight.MasterPassword),
						uint8(3),
					)

					if err == nil {
						log.Println("secret updated", secret)

						ui.InPlainSight.Trigger(events.Event{
							CreatedAt: time.Now(),
							EventType: events.UpdatedSecret,
							Data: map[string]interface{}{
								"secret": secret,
							},
						})

						pages.Navigate("list")
						break
					} else {
						log.Println("error while concealing", err)
					}
				}
			}
		}).
		AddButton("Back", func() {
			pages.Navigate("list")
		})

	grid := tview.NewGrid().
		SetRows(0, 0).
		SetColumns(0, 0, 0)

	flex := tview.NewFlex()
	flex.SetTitle(fmt.Sprintf(" new - inplainsight v%s ", ui.Version)).
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
