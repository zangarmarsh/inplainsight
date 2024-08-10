package secrets

const SecretSeparator byte = '\x02'

type AbstractSecret struct {
	MagicNumber int
	Container   *Container
	deleatable  bool
}

type SecretInterface interface {
	Serialize() string
	UnserializeSecret(serialized string) bool
}

func (s *AbstractSecret) MarkDeleatable() {
	s.deleatable = true
}

func (s *AbstractSecret) IsDeleatable() bool {
	return s.deleatable
}
