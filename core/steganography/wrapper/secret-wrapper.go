package wrapper

import (
	"errors"
	"fmt"
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"github.com/zangarmarsh/inplainsight/core/steganography/header"
	"math"
	"reflect"
	"unsafe"
)

type SecretWrapperInterface interface {
	Len() uint64
	Cap() uint64
	Interweave()
	Unravel()
}

type SecretWrapper struct {
	header   *header.Header
	path     string
	resource any

	data        Secret
	isEncrypted bool
}

type Secret struct {
	encrypted string
	decrypted string
}

func cutYarnChunks(c chan uint8, yarn []uint8, bits int) {
	iterations := int(unsafe.Sizeof(yarn[0])) * 8 / bits
	for _, singleByte := range yarn {
		bitmask := uint8(math.Pow(2, float64(bits)) - 1)

		for i := 0; i < iterations; i++ {
			c <- bitmask & (singleByte >> (8 - (i+1)*bits))
		}
	}

	close(c)
}

func (s *SecretWrapper) craftYarn(secret string) ([]byte, error) {
	if len(secret) == 0 {
		return nil, errors.New("can't add empty secret")
	}

	var buffer []byte

	// Turn each header property into a `byte` sized value and add it to the very beginning of the buffer
	{
		headerData := reflect.ValueOf(*s.header)

		for i := 0; i < headerData.NumField(); i++ {
			buffer = append(buffer, byte(headerData.Field(i).Uint()))
		}
	}

	/**
	 *
	 * Unpacking every rune - hence `int32` data - in four consecutive `byte` sized values
	 *
	 * It will ensure each character to be exactly 4x8 bits and thus grants more consistency
	 * when interweaving/unraveling stuff into media supports.
	 *
	 * The native golang conversion to []byte did not fit here because of how it internally handles
	 * and returns the data, which can be between 1 and 4 bytes long.
	 *
	 */
	{
		const uint8InOneRune = int(unsafe.Sizeof('\x00') / unsafe.Sizeof(uint8(0)))
		for _, singleChar := range secret {
			var unpackedRune [uint8InOneRune]byte

			for shift := 0; shift < uint8InOneRune; shift++ {
				unpackedRune[shift] = byte(singleChar >> (int(unsafe.Sizeof(uint8(0))) * 8 * (uint8InOneRune - shift - 1)))
			}

			buffer = append(buffer, unpackedRune[:]...)
		}
	}

	// ToDo fix this shit
	buffer = append(
		buffer,
		steganography.EndOfMessage,
		steganography.EndOfMessage,
		steganography.EndOfMessage,
		steganography.EndOfMessage,
	)

	return buffer, nil
}

// Todo: fix them up
func (s *SecretWrapper) Interweave(buffer []byte) { fmt.Println("this is fucked up ~~~~~") }
func (s *SecretWrapper) Unravel()                 {}
