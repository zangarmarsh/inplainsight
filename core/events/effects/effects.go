package effects

import (
	"github.com/zangarmarsh/inplainsight/core/events"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/core/utility/lockscreendetector"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"log"
)

// Generic side effects management
func init() {
	var stopCaringAboutLockScreen chan bool

	inplainsight.InPlainSight.AddEventsListener([]events.EventType{events.AppInit, events.UserPreferenceChanged}, func(event events.Event) {

		// Toggle screen lock watcher
		{
			if event.EventType == events.AppInit || (event.Data["pointer"] != nil && event.Data["pointer"] == &inplainsight.InPlainSight.UserPreferences.LogoutOnScreenLock) {
				if stopCaringAboutLockScreen == nil {
					stopCaringAboutLockScreen = make(chan bool)
				}

				if inplainsight.InPlainSight.UserPreferences.LogoutOnScreenLock {
					go func() {
						locked := lockscreendetector.Analyze(&stopCaringAboutLockScreen)

						if <-*locked {
							log.Println("The screen have been locked, logging out...")
							inplainsight.InPlainSight.Logout()
						}
					}()
				} else {
					stopCaringAboutLockScreen <- true
				}
			}
		}
	})

	inplainsight.InPlainSight.AddEventsListener([]events.EventType{events.AppLogout}, func(event events.Event) {
		err := pages.Navigate("register")
		if err != nil {
			log.Println(err)
		}
	})
}
