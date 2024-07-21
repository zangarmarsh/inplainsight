package header

import (
	"errors"
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"reflect"
	"unsafe"
)

type Header struct {
	MagicNumber           uint8
	Version               uint8
	Compression           uint8
	EndOfMessageDelimiter uint8
}

func (h *Header) Bits() int {
	v := reflect.ValueOf(*h)
	return v.NumField() * int(unsafe.Sizeof(uint8(0))) * 8
}

func NewHeader(maximumCompression, endOfMessage uint8) (*Header, error) {
	if maximumCompression < 1 || maximumCompression > 8 {
		return nil, errors.New("compression level must be between 1 and 8")
	}

	header := Header{
		MagicNumber:           steganography.MagicNumber,
		EndOfMessageDelimiter: endOfMessage,
		Version:               steganography.Version,
	}

	return &header, nil
}
