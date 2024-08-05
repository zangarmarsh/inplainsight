package list

import (
	"fmt"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/ui/pages/newsecret"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui/events"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"github.com/zangarmarsh/inplainsight/ui/pages/editsecret"
	"github.com/zangarmarsh/inplainsight/ui/services/logging"
	"github.com/zangarmarsh/inplainsight/ui/widgets"
	"golang.design/x/clipboard"
)

type Page struct {
	pages.GridPage
}

type pageFactory struct{}

func (r pageFactory) GetName() string {
	return "list"
}

var logBox *logging.LogsBox
var filteredSecrets []*inplainsight.Secret
var selectedListItem *int
var searchQuery string

func (r pageFactory) Create() pages.PageInterface {
	var filterResults = func(resultList *tview.List, secrets []*inplainsight.Secret) {
		for path, secret := range secrets {
			log.Printf("%+v in file %+v\n", secret, path)
		}

		lowerCaseSearchQuery := strings.ToLower(strings.TrimLeft(searchQuery, " "))
		resultList.Clear()
		filteredSecrets = nil

		resultList.SetChangedFunc(func(i int, m string, s string, shortcut rune) {
			selectedListItem = &i
		})

		for index, secret := range secrets {
			pasteIntoClipboard := func() {
				log.Println("Copying into clipboard")
				clipboard.Write(clipboard.FmtText, []byte(filteredSecrets[*selectedListItem].Secret))
				logBox.AddLine(fmt.Sprintf("Container '%s' copied into clipboard!", filteredSecrets[*selectedListItem].Secret), logging.Info)
				logBox.AddSeparator()
			}

			if len(searchQuery) == 0 ||
				strings.Contains(strings.ToLower(secret.Title), lowerCaseSearchQuery) ||
				strings.Contains(strings.ToLower(secret.Description), lowerCaseSearchQuery) {
				filteredSecrets = append(
					filteredSecrets,
					secret,
				)

				resultList.InsertItem(
					index,
					secret.Title,
					secret.Description,
					0,
					pasteIntoClipboard,
				)
			}
		}

		if resultList.GetItemCount() == 0 {
			resultList.AddItem("No secrets have been found for the given query.", "", '\x00', nil)
		}
	}

	page := &Page{}
	page.SetName(r.GetName())

	container := tview.NewFlex()
	container.SetBorderPadding(0, 0, 1, 1)
	container.
		SetTitle(fmt.Sprintf(" inplainsight v%s ", inplainsight.Version)).
		SetBorder(true)

	container.SetDirection(tview.FlexRow)

	// Log box
	{
		primitive := tview.NewTextView()
		primitive.
			SetTitle("Logs").
			SetBorder(true)
		primitive.SetDynamicColors(true)
		primitive.SetBorderPadding(0, 0, 2, 2)
		logBox = logging.NewLogsBox(primitive)

		logBox.AddLine("Starting up...", logging.Info)
		logBox.AddLine("Scanning directory...", logging.Info)
		logBox.AddSeparator()
	}

	// Results
	resultBox := tview.NewFlex().SetDirection(tview.FlexRow)
	resultList := tview.NewList()
	resultList.SetBorder(true)
	resultList.SetBorderPadding(0, 0, 2, 2)

	resultBox.AddItem(resultList, 0, 1, false)
	resultBox.SetTitle("Results")

	listOfShortcuts := []string{
		"[orange:italic][ENTER][-] Clipboard",
		"[orange:italic][^N][-] New",
		"[orange:italic][^D][-] Delete",
		"[orange:italic][^E][-] Edit",
		"[orange:italic][^C][-] Quit",
	}

	shortcuts := tview.NewTextView()
	shortcuts.SetDynamicColors(true)
	shortcuts.SetTextAlign(tview.AlignCenter)
	shortcuts.SetText(strings.Join(listOfShortcuts, " "))
	resultBox.AddItem(shortcuts, 1, 0, false)

	// Query box
	queryBox := tview.NewGrid()
	queryBox.SetSize(1, 1, 1, 1)
	queryInput := tview.NewInputField().
		SetPlaceholder("Search for anything")
	queryInput.SetBorderPadding(0, 0, 1, 1)
	queryBox.AddItem(queryInput, 1, 1, 1, 1, 0, 0, true)
	queryBox.SetBorder(true).SetBorderPadding(0, 0, 0, 0)
	queryInput.SetChangedFunc(func(text string) {
		searchQuery = text
		filterResults(resultList, inplainsight.InPlainSight.Secrets)
	})

	queryInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyDown {
			queryInput.Blur()
			resultList.Focus(nil)
		}

		return event
	})

	resultList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyBacktab {
			resultList.Blur()
			queryInput.Focus(nil)
		}

		if event.Key() == tcell.KeyBS || event.Key() == tcell.KeyBackspace2 || (event.Rune() > '\x20' && event.Rune() < '\x7f') {
			logBox.AddLine("Searching for the given key...", logging.Info)
			logBox.AddSeparator()
			resultList.Blur()

			(queryInput.InputHandler())(event, nil)

			return nil
		}

		return event
	})

	container.
		AddItem(queryBox, 4, 0, true).
		AddItem(resultBox, 0, 1, false).
		AddItem(logBox.GetPrimitive(), 11, 0, false)

	copyright := tview.NewTextView()
	copyright.SetDynamicColors(true)
	copyright.SetTextAlign(tview.AlignRight)
	copyright.SetText("[orange:italic] github.com/zangarmarsh/inplainsight")

	container.AddItem(copyright, 1, 1, false)
	page.SetPrimitive(container)

	inplainsight.InPlainSight.AddEventsListener(
		[]events.EventType{events.DiscoveredNewSecret},
		func(event events.Event) {
			resultList.AddItem(
				event.Data["secret"].(*inplainsight.Secret).Title,
				event.Data["secret"].(*inplainsight.Secret).Description,
				0,
				nil,
			)

			filterResults(resultList, inplainsight.InPlainSight.Secrets)
			inplainsight.InPlainSight.App.ForceDraw()

			logLine := event.Data["secret"].(*inplainsight.Secret).Title
			if event.Data["secret"].(*inplainsight.Secret).Description != "" {
				logLine = logLine + " - " + event.Data["secret"].(*inplainsight.Secret).Description
			}

			logBox.AddLine(fmt.Sprintf("Found secret '%s' in file", logLine), logging.Info)
			logBox.AddSeparator()
		})

	inplainsight.InPlainSight.AddEventsListener(
		[]events.EventType{events.UpdatedSecret},
		func(event events.Event) {
			log.Println("event", event)

			filterResults(resultList, inplainsight.InPlainSight.Secrets)
			inplainsight.InPlainSight.App.ForceDraw()
		},
	)

	inplainsight.InPlainSight.AddEventsListener(
		[]events.EventType{events.AddedNewSecret},
		func(event events.Event) {
			resultList.AddItem(event.Data["secret"].(*inplainsight.Secret).Secret, event.Data["secret"].(*inplainsight.Secret).Description, 0, nil)
			filterResults(resultList, inplainsight.InPlainSight.Secrets)
			logBox.AddLine("Added a new secret", logging.Info)
			logBox.AddSeparator()
		})

	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN:
			log.Println("Detected ctrl + n")

			if page := newsecret.Create(); page == nil {
				widgets.ModalError("Generic error")
			} else {
				inplainsight.InPlainSight.Pages.AddAndSwitchToPage(newsecret.GetName(), page.GetPrimitive(), true)
			}

		case tcell.KeyCtrlE:
			if page := editsecret.Create(filteredSecrets[*selectedListItem]); page == nil {
				widgets.ModalError("Generic error")
			} else {
				inplainsight.InPlainSight.Pages.AddAndSwitchToPage(editsecret.GetName(), page.GetPrimitive(), true)
			}

		case tcell.KeyCtrlD:
			widgets.ModalError("Are you sure you want to delete this secret?")
			inplainsight.InPlainSight.App.ForceDraw()

			filteredSecrets[*selectedListItem].Secret = ""
			filteredSecrets[*selectedListItem].Description = ""
			filteredSecrets[*selectedListItem].Title = ""

			// ToDo: find a better way to remove any secret
			err := inplainsight.Conceal(filteredSecrets[*selectedListItem])
			if err != nil {
				return nil
			}

			filteredSecrets = append(filteredSecrets[*selectedListItem:], filteredSecrets[:(*selectedListItem)+1]...)
			filterResults(resultList, inplainsight.InPlainSight.Secrets)

		default:

		}

		return event
	})

	return page
}

func register(records pages.PageFactoryDictionary) pages.PageFactoryInterface {
	records["list"] = pageFactory{}

	return records["list"]
}

var _ = register(pages.PageFactories)
