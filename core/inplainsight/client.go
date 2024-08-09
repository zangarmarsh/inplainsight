package inplainsight

import (
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/events"
	"github.com/zangarmarsh/inplainsight/core/utility/config"
	"time"
)

type InPlainSightClient struct {
	events.EventeableStruct

	App   *tview.Application
	Pages *tview.Pages

	Secrets []*Secret

	Hosts HostsPool

	MasterPassword string
	Path           string

	UserPreferences *config.Config
}

func (c *InPlainSightClient) Logout() {
	c.Secrets = nil
	c.Hosts.Reset()

	c.MasterPassword = ""
	c.Path = ""

	c.UserPreferences = nil

	c.Trigger(events.Event{
		CreatedAt: time.Now(),
		EventType: events.AppLogout,
		Data:      map[string]interface{}{},
	})
}
