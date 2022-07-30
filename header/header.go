package header

import (
	"reflect"
	"unsafe"
)

type Header struct {
	MagicNumber           uint8
	Version               uint8
	Compression           uint8
	EndOfCharDelimiter    uint8
	EndOfMessageDelimiter uint8
}

func (h *Header) Size() int {
	v := reflect.ValueOf( *h )
	return v.NumField() * int(unsafe.Sizeof( uint8(0) )) * 8
}

