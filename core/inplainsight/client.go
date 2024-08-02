package inplainsight

import (
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui/events"
)

type InPlainSightClient struct {
	events.EventeableStruct

	App   *tview.Application
	Pages *tview.Pages

	Secrets map[string]*Secret
	Hosts   PoolOfHosts

	MasterPassword string
	Path           string
}
