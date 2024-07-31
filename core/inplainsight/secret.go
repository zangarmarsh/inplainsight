package inplainsight

import (
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"log"
	"strings"
)

type Secret struct {
	Title       string
	Description string
	Secret      string
	FilePath    string

	Host steganography.SecretInterface
}

func (s *Secret) Serialize() string {
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
