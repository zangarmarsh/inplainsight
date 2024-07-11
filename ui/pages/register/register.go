package register

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

type Page struct {
	pages.GridPage
}

type pageFactory struct{}

func (r pageFactory) GetName() string {
	return "register"
}

func (r pageFactory) Create() pages.PageInterface {
	page := Page{}
	page.SetName(r.GetName())

	form := tview.NewForm()

	form.
		SetBorder(false)

	form.
		AddPasswordField("Master Password", "", 0, '*', nil).
		AddInputField("Pool path", "~/Pictures/passwords/", 0, nil, nil).
		SetButtonsAlign(tview.AlignCenter).
		AddButton("Register", func() {
			password := form.GetFormItemByLabel("Master Password").(*tview.InputField).GetText()
			path := form.GetFormItemByLabel("Pool path").(*tview.InputField).GetText()

			if path[0] == '~' {
				homeDir, err := os.UserHomeDir()
				if err == nil {
					path = fmt.Sprintf("%s/%s", homeDir, strings.TrimLeft(path[1:], "/\\"))
				}
			}

			files, err := os.ReadDir(path)
			if err != nil {
				log.Println(err)
				widgets.ModalError(err.Error())
				ui.InPlainSight.App.ForceDraw()
				return
			}

			if len(files) == 0 {
				log.Println("empty directory")
			}

			ui.InPlainSight.Path = path
			ui.InPlainSight.MasterPassword = password

			err = pages.Navigate("list")

			s := steganography.Steganography{}
			var eligibleFiles []os.DirEntry

			for _, file := range files {
				go revealSecret(file, eligibleFiles, path, err, s, password)
			}

			// if err != nil {
			// 	widgets.ModalError("Generic error")
			// }
		}).
		AddButton("Quit (CTRL + C)", func() {
			ui.InPlainSight.App.Stop()
		})

	grid := tview.NewGrid().
		SetRows(0, 0).
		SetColumns(0, 0, 0)

	flex := tview.NewFlex()
	flex.SetTitle(fmt.Sprintf(" register - inplainsight v%s ", ui.Version)).
		SetBorder(true)

	flex.SetDirection(tview.FlexRow)
	flex.SetBorderPadding(2, 2, 2, 2)

	text := tview.NewTextView().
		SetText(
			"Greetings, stranger. To get it started you'll just have to set a strong password.\n\n" +
				"[red::b]Beware:[-:-:-] don't forget it or you won't be able to access to your secrets never again.",
		).
		SetWordWrap(true).
		SetDynamicColors(true)

	flex.AddItem(text, 0, 2, false)
	flex.AddItem(form, 0, 1, true)

	grid.
		AddItem(flex, 0, 1, 1, 1, 30, 50, true).
		AddItem(flex, 0, 0, 3, 3, 0, 0, true)

	grid.SetBorderPadding(2, 0, 0, 0)

	page.SetPrimitive(grid)

	return &page
}

func revealSecret(file os.DirEntry, eligibleFiles []os.DirEntry, path string, err error, s steganography.Steganography, password string) {
	// ToDo: refactor this entire block keeping in mind that there's a whole lotta of file extensions

	if !file.IsDir() && strings.Contains(file.Name(), ".png") {
		eligibleFiles = append(eligibleFiles, file)
		log.Println("found eligibile file " + file.Name())

		filePath := fmt.Sprintf("%s/%s", strings.TrimRight(path, "/\\"), file.Name())
		var revealed string
		revealed, err = s.Reveal(filePath, []byte(password))
		log.Println(fmt.Sprintf("master password used to reveal %#v", password))

		if err == nil {
			log.Println(fmt.Sprintf("found secret %#v", revealed))

			secret := &ui.Secret{}
			secret.Unserialize(revealed)
			secret.FilePath = filePath

			ui.InPlainSight.InvolvedFiles = append(ui.InPlainSight.InvolvedFiles, file.Name())
			ui.InPlainSight.Secrets = append(ui.InPlainSight.Secrets, secret)

			ui.InPlainSight.Trigger(events.Event{
				CreatedAt: time.Now(),
				EventType: events.DiscoveredNewSecret,
				Data: map[string]interface{}{
					"secret": *secret,
				},
			})
		} else {
			log.Println("theres no secret in here")
		}
	}
}

func register(records pages.PageFactoryDictionary) pages.PageFactoryInterface {
	records["register"] = pageFactory{}
	return records["register"]
}

var _ = register(pages.PageFactories)
