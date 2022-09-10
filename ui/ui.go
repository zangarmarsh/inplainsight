package ui

import (
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui/events"
	"golang.design/x/clipboard"
	"log"
	"strings"
)

type Secret struct {
	Title        string
	Description  string
	Secret       string
	FilePath     string
}

const secretSeparator = "\x02"

func (s *Secret ) Serialize() string  {
	return s.Title + secretSeparator + s.Description + secretSeparator + s.Secret
}

func (s *Secret) Unserialize(serialized string) {
	fields := strings.Split(serialized, secretSeparator)

	if len(fields) == 3 {
		s.Title = fields[0]
		s.Description = fields[1]
		s.Secret = fields[2]
	} else {
		log.Printf("Cannot unserialize serialized secret %#v\n", serialized)
	}
}

type InPlainSightClient struct {
	events.EventeableStruct

	App            *tview.Application
	Pages          *tview.Pages

	InvolvedFiles  []string
	Secrets        []*Secret

	MasterPassword string
	Path           string
}

var InPlainSight = &InPlainSightClient{}

const (
	Version = "1.0.0"
)

func Bootstrap() {
	err := clipboard.Init()
	if err != nil {
		return
	}

	InPlainSight.App   = tview.NewApplication()
	InPlainSight.Pages = tview.NewPages()

	InPlainSight.App.SetRoot(InPlainSight.Pages, true).EnableMouse(true)
}
