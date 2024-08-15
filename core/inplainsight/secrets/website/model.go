package website

import (
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"log"
	"strings"
)

const magicNumber secrets.MagicNumber = 0x02

type WebsiteCredential struct {
	secrets.AbstractSecret

	website  string
	note     string
	account  string
	password string
}

func init() {
	secrets.SecretsModelRegister[magicNumber] = func(serialized string) secrets.SecretInterface {
		return (&WebsiteCredential{}).Unserialize(serialized)
	}
}

func (s *WebsiteCredential) Serialize() string {
	if s.Deleatable {
		return ""
	}

	return string(s.GetMagicNumber()) +
		s.website + string(secrets.SecretSeparator) +
		s.note + string(secrets.SecretSeparator) +
		s.account + string(secrets.SecretSeparator) +
		s.password
}

func (s *WebsiteCredential) Unserialize(serialized string) secrets.SecretInterface {
	fields := strings.Split(serialized, string(secrets.SecretSeparator))

	if len(fields) == 4 {
		s.website = fields[0]
		s.note = fields[1]
		s.account = fields[2]
		s.password = fields[3]

		return s
	} else {
		log.Printf("Cannot unserialize serialized password %#v\n", serialized)
		return nil
	}
}

func (s *WebsiteCredential) Filter(query string) bool {
	query = strings.ToLower(strings.Trim(query, " "))

	return strings.Contains(strings.ToLower(s.website), query) ||
		strings.Contains(strings.ToLower(s.note), query) ||
		strings.Contains(strings.ToLower(s.account), query)
}

func (s *WebsiteCredential) GetMagicNumber() secrets.MagicNumber {
	return magicNumber
}

func (s *WebsiteCredential) SetTitle(website string) {
	s.website = website

}
func (s *WebsiteCredential) GetTitle() string {
	return s.website
}

func (s *WebsiteCredential) SetDescription(note string) {
	s.note = note
}

func (s *WebsiteCredential) GetDescription() string {
	return s.note
}

func (s *WebsiteCredential) GetIcon() rune {
	return '🌐'
}
