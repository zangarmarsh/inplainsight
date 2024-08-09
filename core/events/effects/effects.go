package effects

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/events"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/core/utility/lockscreendetector"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"log"
	"time"
)

// Generic side effects management
func init() {
	var stopCaringAboutLockScreen chan bool
	var interactionsChannel chan bool

	inplainsight.InPlainSight.AddEventsListener([]events.EventType{events.AppLogout}, func(event events.Event) {
		inplainsight.InPlainSight.Pages.RemovePage("list")
		err := pages.Navigate("register")
		inplainsight.InPlainSight.App.Draw()
		if err != nil {
			log.Println(err)
		}

	})

	inplainsight.InPlainSight.AddEventsListener([]events.EventType{events.AppLogout}, func(event events.Event) {
		log.Println("cleaning up before logging out")

		// Resetting mouse and keyboard event capture prevents a channel deadlock
		inplainsight.InPlainSight.App.SetInputCapture(nil)
		inplainsight.InPlainSight.App.SetMouseCapture(nil)

		if stopCaringAboutLockScreen != nil {
			stopCaringAboutLockScreen <- true
			close(stopCaringAboutLockScreen)
		}

		if interactionsChannel != nil {
			close(interactionsChannel)
		}
	})

	inplainsight.InPlainSight.AddEventsListener([]events.EventType{events.AppInit, events.UserPreferenceChanged}, func(event events.Event) {
		// Toggle screen lock watcher
		{
			if event.EventType == events.AppInit || (event.Data["pointer"] != nil && event.Data["pointer"] == &inplainsight.InPlainSight.UserPreferences.LogoutOnScreenLock) {
				if inplainsight.InPlainSight.UserPreferences.LogoutOnScreenLock {
					go func(ch *chan bool) {
						if *ch == nil {
							*ch = make(chan bool)
						}

						locked := lockscreendetector.Analyze(ch)

						if <-*locked {
							log.Println("The screen have been locked, logging out...")
							inplainsight.InPlainSight.Logout()
						}
					}(&stopCaringAboutLockScreen)
				} else if stopCaringAboutLockScreen != nil {
					stopCaringAboutLockScreen <- true
				}
			}
		}

		// Toggle AFK timeout
		{
			// Todo implement countdown using https://pkg.go.dev/github.com/rivo/tview#Application.SetInputCapture
			//  and https://pkg.go.dev/github.com/rivo/tview#Application.SetMouseCapture

			if event.EventType == events.AppInit || (event.Data["pointer"] != nil &&
				event.Data["pointer"] == &inplainsight.InPlainSight.UserPreferences.AFKTimeout) {
				var afkTimer *time.Timer

				if inplainsight.InPlainSight.UserPreferences.AFKTimeout != 0 {

					if inplainsight.InPlainSight.UserPreferences.AFKTimeout < 5 {
						inplainsight.InPlainSight.UserPreferences.AFKTimeout = 5
					}

					if afkTimer != nil {
						afkTimer.Stop()
						afkTimer.Reset(time.Minute * time.Duration(inplainsight.InPlainSight.UserPreferences.AFKTimeout))
					} else {
						afkTimer = time.NewTimer(time.Minute * time.Duration(inplainsight.InPlainSight.UserPreferences.AFKTimeout))
					}

					go func() {
						for {
							select {
							case _, ok := <-afkTimer.C:
								if !ok {
									return
								}

								log.Println("dropped after afk timeout")

								inplainsight.InPlainSight.Logout()
								return
							case _, ok := <-interactionsChannel:
								if !ok {
									log.Println("closed interactions channel")
									return
								}

								afkTimer.Reset(time.Minute * time.Duration(inplainsight.InPlainSight.UserPreferences.AFKTimeout))
							default:
								time.Sleep(time.Millisecond * 10)
							}
						}
					}()

					go func() {
						inplainsight.InPlainSight.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
							if interactionsChannel == nil {
								interactionsChannel = make(chan bool)
							}

							interactionsChannel <- true

							return event
						})

						inplainsight.InPlainSight.App.SetMouseCapture(
							func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
								if interactionsChannel == nil {
									interactionsChannel = make(chan bool)
								}

								interactionsChannel <- true

								return event, action
							})
					}()

				} else {
					if afkTimer != nil {
						afkTimer.Stop()
					}

					if interactionsChannel != nil {
						close(interactionsChannel)
					}
				}
			}
		}
	})
}
