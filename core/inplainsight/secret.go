package inplainsight

import (
	"log"
	"strings"
)

// ToDo: add magic number
type Secret struct {
	Title       string
	Description string
	Secret      string

	Container  *SecretsContainer
	deleatable bool
}

func (s *Secret) Serialize() string {
	if s.deleatable {
		return ""
	}

	return s.Title + secretSeparator + s.Description + secretSeparator + s.Secret
}

func UnserializeSecret(serialized string) *Secret {
	secret := Secret{}

	fields := strings.Split(serialized, secretSeparator)

	if len(fields) == 3 {
		secret.Title = fields[0]
		secret.Description = fields[1]
		secret.Secret = fields[2]

		return &secret
	} else {
		log.Printf("Cannot unserialize serialized secret %#v\n", serialized)
		return nil
	}
}

func (s *Secret) MarkDeleatable() {
	s.deleatable = true
}
