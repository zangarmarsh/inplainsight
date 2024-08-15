package simple

import (
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"log"
	"strings"
)

const magicNumber secrets.MagicNumber = 0x01

type SimpleSecret struct {
	secrets.AbstractSecret

	title       string
	description string
	secret      string
}

func init() {
	secrets.SecretsModelRegister[magicNumber] = func(serialized string) secrets.SecretInterface {
		return (&SimpleSecret{}).Unserialize(serialized)
	}
}

func (s *SimpleSecret) Serialize() string {
	if s.Deleatable {
		return ""
	}

	return string(s.GetMagicNumber()) + s.title + string(secrets.SecretSeparator) + s.description + string(secrets.SecretSeparator) + s.secret
}

func (s *SimpleSecret) Unserialize(serialized string) secrets.SecretInterface {
	fields := strings.Split(serialized, string(secrets.SecretSeparator))

	if len(fields) == 3 {
		s.title = fields[0]
		s.description = fields[1]
		s.secret = fields[2]

		return s
	} else {
		log.Printf("Cannot unserialize serialized secret %#v\n", serialized)
		return nil
	}
}

func (s *SimpleSecret) Filter(query string) bool {
	query = strings.ToLower(strings.Trim(query, " "))

	return strings.Contains(strings.ToLower(s.GetTitle()), query) ||
		strings.Contains(strings.ToLower(s.GetDescription()), query)
}

func (s *SimpleSecret) GetMagicNumber() secrets.MagicNumber {
	return magicNumber
}

func (s *SimpleSecret) SetTitle(title string) {
	s.title = title

}
func (s *SimpleSecret) GetTitle() string {
	return s.title
}

func (s *SimpleSecret) SetDescription(description string) {
	s.description = description

}
func (s *SimpleSecret) GetDescription() string {
	return s.description
}

func (s *SimpleSecret) SetSecret(secret string) {
	s.secret = secret

}
func (s *SimpleSecret) GetSecret() string {
	return s.secret
}
