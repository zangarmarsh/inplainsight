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

// Given a rune it splits it up to 4 uint8 values, then assigns them to the header
// Returns `true` if the magic number is matched
func (h *Header) Set(r rune) (bool, error) {
	bitmask := ^rune(0)

	if h == nil {
		h = &Header{}
	}

	for counter, fieldPtr := range []*uint8{&h.MagicNumber, &h.Version, &h.Compression, &h.EndOfMessageDelimiter} {
		*fieldPtr = uint8(bitmask & (r >> (int(unsafe.Sizeof(r))*8 - counter*8 - int(unsafe.Sizeof(uint8(0))*8))))
	}

	return h.MagicNumber == steganography.MagicNumber, nil
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
