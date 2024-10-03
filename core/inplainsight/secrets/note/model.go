package note

import (
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"log"
	"strconv"
	"strings"
)

const magicNumber secrets.MagicNumber = 0x04

type Note struct {
	secrets.AbstractSecret

	title    string
	isHidden bool
	note     string
}

func init() {
	secrets.SecretsModelRegister[magicNumber] = func(serialized string) secrets.SecretInterface {
		return (&Note{}).Unserialize(serialized)
	}
}

func (s *Note) Serialize() string {
	if s.Deleatable {
		return ""
	}

	return string(s.GetMagicNumber()) +
		s.title + string(secrets.SecretSeparator) +
		strconv.FormatBool(s.isHidden) + string(secrets.SecretSeparator) +
		s.note
}

func (s *Note) Unserialize(serialized string) secrets.SecretInterface {
	fields := strings.Split(serialized, string(secrets.SecretSeparator))

	if len(fields) == 3 {
		s.title = fields[0]
		s.isHidden = fields[1] == "true"
		s.note = fields[2]

		return s
	} else {
		log.Printf("Cannot unserialize secret %#v\n", serialized)
		return nil
	}
}

func (s *Note) Filter(query string) bool {
	return strings.Contains(s.title, query)
}

func (s *Note) GetMagicNumber() secrets.MagicNumber {
	return magicNumber
}

func (s *Note) SetTitle(website string) {

}
func (s *Note) GetTitle() string {
	return s.title
}

func (s *Note) SetDescription(note string) {
	return
}

func (s *Note) GetDescription() string {
	return s.note
}

func (s *Note) GetIcon() rune {
	return '📝'
}
