package secrets

import (
	"log"
	"strings"
)

const MagicNumber = 0x01

type SimpleSecret struct {
	AbstractSecret

	Title       string
	Description string
	Secret      string
}

func (s *SimpleSecret) Serialize() string {
	if s.deleatable {
		return ""
	}

	return s.Title + string(SecretSeparator) + s.Description + string(SecretSeparator) + s.Secret
}

func (s *SimpleSecret) UnserializeSecret(serialized string) *SimpleSecret {
	fields := strings.Split(serialized, string(SecretSeparator))

	if len(fields) == 3 {
		s.Title = fields[0]
		s.Description = fields[1]
		s.Secret = fields[2]

		return s
	} else {
		log.Printf("Cannot unserialize serialized secret %#v\n", serialized)
		return nil
	}
}
