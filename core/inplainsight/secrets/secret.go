package secrets

import (
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

const SecretSeparator byte = '\x02'

type MagicNumber byte

type SecretInterface interface {
	Serialize() string
	Unserialize(serialized string) SecretInterface

	SetHeader(header *Header)
	GetHeader() *Header

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

	AssignRandomID()
	GetID() string
}

type AbstractSecret struct {
	id         string
	Deleatable bool
	container  *Container

	header *Header
}

func (s *AbstractSecret) AssignRandomID() {
	// ToDo handle error cases
	generatedUUID, _ := uuid.NewRandom()
	s.id = generatedUUID.String()
}

func (s *AbstractSecret) GetID() string {
	return s.id
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

func (s *AbstractSecret) SetHeader(h *Header) {
	s.header = h
}

func (s *AbstractSecret) GetHeader() *Header {
	return s.header
}

func LinkSecretAndContainer(secret SecretInterface, container *Container) {
	secret.SetContainer(container)
	(*container).Add(secret)
}
