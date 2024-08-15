package file

import (
	"encoding/base64"
	"github.com/zangarmarsh/inplainsight/core/inplainsight/secrets"
	"log"
	"strings"
)

const magicNumber secrets.MagicNumber = 0x05

type File struct {
	secrets.AbstractSecret

	title    string
	note     string
	fileName string
	content  string
}

func init() {
	secrets.SecretsModelRegister[magicNumber] = func(serialized string) secrets.SecretInterface {
		return (&File{}).Unserialize(serialized)
	}
}

func (s *File) Serialize() string {
	if s.Deleatable {
		return ""
	}

	serializedContent := base64.StdEncoding.EncodeToString([]byte(s.content))

	return string(s.GetMagicNumber()) +
		s.title + string(secrets.SecretSeparator) +
		s.note + string(secrets.SecretSeparator) +
		s.fileName + string(secrets.SecretSeparator) +
		serializedContent
}

func (s *File) Unserialize(serialized string) secrets.SecretInterface {
	fields := strings.Split(serialized, string(secrets.SecretSeparator))

	if len(fields) == 4 {
		s.title = fields[0]
		s.note = fields[1]
		s.fileName = fields[2]
		s.content = fields[3]

		decodeContent, err := base64.StdEncoding.DecodeString(s.content)
		if len(decodeContent) > 0 && err == nil {
			s.content = string(decodeContent)
			return s
		} else {
			log.Printf("Cannot base64 decode string '%s': %v", serialized, err)
		}
	} else {
		log.Printf("Cannot unserialize secret %#v\n", serialized)
	}

	return nil
}

func (s *File) Filter(query string) bool {
	return strings.Contains(s.title, query) || strings.Contains(s.note, query) || strings.Contains(s.fileName, query)
}

func (s *File) GetMagicNumber() secrets.MagicNumber {
	return magicNumber
}

func (s *File) SetTitle(title string) { s.title = title }
func (s *File) GetTitle() string {
	return s.title
}

func (s *File) SetDescription(note string) {
	s.note = note
}

func (s *File) GetDescription() string {
	return s.note
}

func (s *File) GetIcon() rune {
	return '📝'
}
