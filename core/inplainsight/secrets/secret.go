package secrets

import (
	"github.com/rivo/tview"
)

const SecretSeparator byte = '\x02'

type MagicNumber byte

type SecretInterface interface {
	Serialize() string
	Unserialize(serialized string) SecretInterface
	MarkDeleatable()
	IsDeleatable() bool

	SetContainer(*Container)
	GetContainer() *Container

	GetTitle() string
	SetTitle(title string)

	GetDescription() string
	SetDescription(description string)

	GetForm() *tview.Form
	Filter(query string) bool

	DoAction()
	GetIcon() rune
}

type AbstractSecret struct {
	Deleatable bool
	container  *Container
}

func (s *AbstractSecret) MarkDeleatable() {
	s.Deleatable = true
}

func (s *AbstractSecret) IsDeleatable() bool {
	return s.Deleatable
}

func (s *AbstractSecret) SetContainer(container *Container) {
	s.container = container
}

func (s *AbstractSecret) GetContainer() *Container {
	return s.container
}

// Generic icon
func (s *AbstractSecret) GetIcon() rune {
	return '👤'
}
