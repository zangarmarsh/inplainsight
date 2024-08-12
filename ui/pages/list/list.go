package list

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/events"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"github.com/zangarmarsh/inplainsight/ui/pages/editsecret"
	"github.com/zangarmarsh/inplainsight/ui/pages/newsecret"
	"github.com/zangarmarsh/inplainsight/ui/services/logging"
	"github.com/zangarmarsh/inplainsight/ui/widgets"
	"log"
	"strconv"
	"strings"
	"time"
)

type Page struct {
	pages.GridPage
}

type pageFactory struct{}

func (r pageFactory) GetName() string {
	return "list"
}

var logBox *logging.LogsBox
var filteredSecrets []secrets.SecretInterface
var selectedListItem *int
var searchQuery string

func (r pageFactory) Create() pages.PageInterface {
	// Todo find a smarter way to filter the results
	var filterResults = func(resultList *tview.List, secrets []secrets.SecretInterface) {
		resultList.Clear()
		filteredSecrets = nil

		resultList.SetChangedFunc(func(i int, m string, s string, shortcut rune) {
			selectedListItem = &i
		})

		for index, secret := range secrets {
			if len(searchQuery) == 0 || secret.Filter(searchQuery) {
				filteredSecrets = append(
					filteredSecrets,
					secret,
				)

				resultList.InsertItem(
					index,
					fmt.Sprintf("%s%s%s%s", strings.Repeat(" ", 1), string(secret.GetIcon()), strings.Repeat(" ", 1), secret.GetTitle()),
					fmt.Sprintf("%s%s", strings.Repeat(" ", 4), secret.GetDescription()),
					0,
					secret.DoAction,
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

	footer := tview.NewFlex()

	// Log box
	{
		primitive := tview.NewTextView()
		primitive.
			SetTitle("Logs").
			SetTitleAlign(tview.AlignLeft).
			SetBorder(true)
		primitive.SetDynamicColors(true)
		primitive.SetBorderPadding(0, 0, 2, 2)
		logBox = logging.NewLogsBox(primitive)

		logBox.AddLine("Starting up...", logging.Info)
		logBox.AddLine("Scanning directory...", logging.Info)
		logBox.AddSeparator()

		footer.AddItem(primitive, 0, 1, false)
	}

	// User preferences form box
	{
		userPreferenceForm := tview.NewForm()

		userPreferenceForm.
			SetBorder(true).
			SetBorderPadding(0, 0, 2, 2).
			SetTitleAlign(tview.AlignRight).
			SetTitle("Preferences")

		// AFK Timeout preference management
		var afkTimeoutInput *tview.InputField
		{
			afkTimeoutInput = tview.NewInputField()
			afkTimeoutInput.
				SetLabel("AFK Timeout (minutes)").
				SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
					return lastChar >= '0' && lastChar <= '9'
				}).
				SetBlurFunc(func() {
					inputValue := afkTimeoutInput.GetText()
					if inputValue != "" {
						inputValue, err := strconv.Atoi(inputValue)

						if err != nil {
							logBox.AddLine("There was an error saving AFK Timeout", logging.Warning)
						} else {
							inplainsight.InPlainSight.UserPreferences.AFKTimeout = inputValue
							err = inplainsight.InPlainSight.UserPreferences.Save()
							if err != nil {
								logBox.AddLine("There was an error saving AFK Timeout", logging.Warning)
							} else {
								inplainsight.InPlainSight.Trigger(events.Event{
									CreatedAt: time.Now(),
									EventType: events.UserPreferenceChanged,
									Data: map[string]interface{}{
										"pointer": &inplainsight.InPlainSight.UserPreferences.AFKTimeout,
									},
								})
							}
						}
					}
				})

			userPreferenceForm.AddFormItem(afkTimeoutInput)
		}

		// Logout on screen lock user preference management
		var logoutOnScreenLockCheckbox *tview.Checkbox
		{
			logoutOnScreenLockCheckbox = tview.NewCheckbox()
			logoutOnScreenLockCheckbox.
				SetLabel("Logout on screen lock").
				SetChangedFunc(func(checked bool) {
					inplainsight.InPlainSight.UserPreferences.LogoutOnScreenLock = checked
					err := inplainsight.InPlainSight.UserPreferences.Save()

					if err != nil {
						logBox.AddLine("There was an error saving Logout on screen lock", logging.Warning)
					} else {
						inplainsight.InPlainSight.Trigger(events.Event{
							CreatedAt: time.Now(),
							EventType: events.UserPreferenceChanged,
							Data: map[string]interface{}{
								"pointer": &inplainsight.InPlainSight.UserPreferences.LogoutOnScreenLock,
							},
						})
					}
				})

			userPreferenceForm.AddFormItem(logoutOnScreenLockCheckbox)
		}

		// This delayed initialization might prevent casual null pointer deference which would occur if the current page is
		// rendered before the user login event
		if inplainsight.InPlainSight.UserPreferences != nil {
			afkTimeoutInput.SetText(strconv.Itoa(inplainsight.InPlainSight.UserPreferences.AFKTimeout))
			logoutOnScreenLockCheckbox.SetChecked(inplainsight.InPlainSight.UserPreferences.LogoutOnScreenLock)
		}

		footer.AddItem(userPreferenceForm, 50, 1, false)
	}

	// Results
	resultBox := tview.NewFlex().SetDirection(tview.FlexRow)

	resultList := tview.NewList()
	resultList.
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2)

	resultList.SetHighlightFullLine(false)

	resultBox.AddItem(resultList, 0, 1, false)
	resultBox.SetTitle("Results")

	// Shortcuts
	{
		listOfShortcuts := []string{
			"[orange:italic][ Space ][-] Action",
			"[orange:italic][ ^N ][-] New",
			"[orange:italic][ ^D ][-] Delete",
			"[orange:italic][ ^E ][-] Edit",
			"[orange:italic][ ^C ][-] Quit",
		}

		shortcuts := tview.NewTextView()
		shortcuts.SetDynamicColors(true)
		shortcuts.SetTextAlign(tview.AlignCenter)
		shortcuts.SetText(strings.Join(listOfShortcuts, " "))
		resultBox.AddItem(shortcuts, 1, 1, false)
	}

	// Currently open source label
	{
		sourceOfDataLabel := tview.NewTextView().
			SetDynamicColors(true).
			SetTextAlign(tview.AlignRight)

		resultBox.AddItem(sourceOfDataLabel, 1, 1, false)

		inplainsight.InPlainSight.AddEventsListener([]events.EventType{events.AppInit}, func(event events.Event) {
			sourceOfDataLabel.SetText(inplainsight.InPlainSight.Path)
		})
	}

	// Query box
	queryBox := tview.NewGrid()
	queryBox.SetSize(1, 1, 1, 1)
	queryInput := tview.NewInputField().
		SetPlaceholder("Search...").
		SetPlaceholderTextColor(tcell.ColorWhite)

	queryInput.SetBorderPadding(0, 0, 1, 1)
	queryBox.AddItem(queryInput, 1, 1, 1, 1, 0, 0, true)
	queryBox.SetBorder(true).SetBorderPadding(0, 0, 0, 0)
	queryInput.SetChangedFunc(func(text string) {
		searchQuery = text
		filterResults(resultList, inplainsight.InPlainSight.Secrets)
	})

	queryInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown {
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
			inplainsight.InPlainSight.App.SetFocus(queryInput)
			resultList.Blur()

			(queryInput.InputHandler())(event, nil)

			return nil
		}

		return event
	})

	container.
		AddItem(queryBox, 4, 0, true).
		AddItem(resultBox, 0, 1, false).
		AddItem(footer, 11, 0, false)

	copyright := tview.NewTextView()
	copyright.SetDynamicColors(true)
	copyright.SetTextAlign(tview.AlignRight)
	copyright.SetText("[orange:italic] github.com/zangarmarsh/inplainsight")

	container.AddItem(copyright, 1, 1, false)
	page.SetPrimitive(container)

	inplainsight.InPlainSight.AddEventsListener(
		[]events.EventType{events.SecretDiscovered},
		func(event events.Event) {
			resultList.AddItem(
				event.Data["secret"].(secrets.SecretInterface).GetTitle(),
				event.Data["secret"].(secrets.SecretInterface).GetDescription(),
				0,
				nil,
			)

			filterResults(resultList, inplainsight.InPlainSight.Secrets)
			inplainsight.InPlainSight.App.ForceDraw()

			logLine := event.Data["secret"].(secrets.SecretInterface).GetTitle()
			if event.Data["secret"].(secrets.SecretInterface).GetDescription() != "" {
				logLine = logLine + " - " + event.Data["secret"].(secrets.SecretInterface).GetDescription()
			}

			// logBox.AddLine(fmt.Sprintf("Found secret '%s' in file", logLine), logging.Info)
			// logBox.AddSeparator()
		})

	inplainsight.InPlainSight.AddEventsListener(
		[]events.EventType{events.SecretUpdated},
		func(event events.Event) {
			filterResults(resultList, inplainsight.InPlainSight.Secrets)
			inplainsight.InPlainSight.App.ForceDraw()
		},
	)

	inplainsight.InPlainSight.AddEventsListener(
		[]events.EventType{events.SecretAdded},
		func(event events.Event) {
			resultList.AddItem(event.Data["secret"].(secrets.SecretInterface).GetTitle(), event.Data["secret"].(secrets.SecretInterface).GetDescription(), 0, nil)
			filterResults(resultList, inplainsight.InPlainSight.Secrets)
			logBox.AddLine("Added a new secret", logging.Info)
			logBox.AddSeparator()
		})

	inplainsight.InPlainSight.AddEventsListener(
		[]events.EventType{events.UserPreferenceChanged},
		func(event events.Event) {
			logBox.AddLine("User preference changed", logging.Info)
			logBox.AddSeparator()
		},
	)

	inplainsight.InPlainSight.AddEventsListener(
		[]events.EventType{events.Navigation},
		func(event events.Event) {
			log.Println("caught navigation event", event.Data["slug"])
			if event.Data["slug"] == "list" {
				log.Println("setting focus into query input..")
				// inplainsight.InPlainSight.App.SetFocus(queryInput)
			} else {
				queryInput.Blur()
				resultList.Blur()
			}
		},
	)

	inplainsight.InPlainSight.AddEventsListener(
		[]events.EventType{events.SecretCopiedIntoClipboard},
		func(_ events.Event) {
			logBox.AddLine("Secret copied into clipboard!", logging.Info)
			logBox.AddSeparator()
		},
	)

	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN:
			if page := newsecret.Create(); page == nil {
				widgets.ModalAlert("Generic error", nil)
			} else {
				pages.Navigate(page)
			}

		case tcell.KeyCtrlE:
			if page := editsecret.Create(filteredSecrets[*selectedListItem]); page == nil {
				widgets.ModalAlert("Generic error", nil)
			} else {
				pages.Navigate(page)
			}

		case tcell.KeyCtrlD:
			widgets.ModalAlert("Are you sure you want to delete this secret?", func() {
				filteredSecrets[*selectedListItem].MarkDeleatable()
				err := inplainsight.Conceal(filteredSecrets[*selectedListItem])
				if err != nil {
					widgets.ModalAlert("There was an error deleting the secret, please retry", func() {
						pages.GoBack()
					})

					return
				}

				filteredSecrets = append(filteredSecrets[*selectedListItem:], filteredSecrets[:(*selectedListItem)+1]...)
				filterResults(resultList, inplainsight.InPlainSight.Secrets)
			})
			inplainsight.InPlainSight.App.ForceDraw()

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
