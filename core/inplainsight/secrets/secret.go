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

	GetSecret() string
	SetSecret(secret string)

	// Todo callingPage is needed since we need to know which page should we move on after the form has been submitted/canceled
	//      find a smart way to remove it and go back automatically
	GetForm() *tview.Form
	Filter(query string) bool

	DoAction()
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
