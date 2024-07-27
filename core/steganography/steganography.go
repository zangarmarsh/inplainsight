package steganography

import (
	"errors"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"reflect"
	"unsafe"
)

const Version uint8 = '\x01'
const MagicNumber uint8 = '\x78'
const EndOfMessage uint8 = '\x00'

type SecretInterface interface {
	Len() uint64
	Cap() uint64

	Interweave(secret string) error
	Unravel(path string) error

	Data() *SecretData
}

type Secret struct {
	Header   *Header
	Path     string
	resource any

	data        SecretData
	isEncrypted bool
}

type SecretData struct {
	Encrypted string
	Decrypted string
}

// Taken an uint8 array in input, split it into chunks of bits `bits` long
func CutYarnChunks(c chan uint8, yarn []uint8, bits int) {
	iterations := int(math.Ceil(float64(unsafe.Sizeof(yarn[0])) * float64(8) / float64(bits)))
	genericBitmask := uint8(math.Pow(2, float64(bits)) - 1)

	for _, singleByte := range yarn {
		for i := 0; i < iterations; i++ {
			offset := 8 - (i+1)*bits
			bitmask := genericBitmask

			// Adapt the bitmask if bits * channels is not a multiplier of 32
			if offset < 0 {
				bitmask >>= int(math.Abs(float64(offset)))
				offset = 0
			}

			c <- bitmask & (singleByte >> offset)
		}
	}

	close(c)
}

// Given a secret, returns an array of `byte` containing the Header in the first place
// and then the segmented secret
func (s *Secret) CraftYarn(secret string) ([]byte, error) {
	if len(secret) == 0 {
		return nil, errors.New("can't add empty secret")
	}

	var buffer []byte

	// Turn each Header property into a `byte` sized value and add it to the very beginning of the buffer
	{
		headerData := reflect.ValueOf(*s.Header)

		for i := 0; i < headerData.NumField(); i++ {
			buffer = append(buffer, byte(headerData.Field(i).Uint()))
		}
	}

	secret += string(EndOfMessage)

	/**
	 *
	 * Unpacking every rune - hence `int32` Data - in four consecutive `byte` sized values
	 *
	 * It will ensure each character to be exactly 4x8 bits and thus grants more consistency
	 * when interweaving/unraveling stuff into media supports.
	 *
	 * The native golang conversion to []byte did not fit here because of how it internally handles
	 * and returns the Data, which can be between 1 and 4 bytes long.
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

	return buffer, nil
}

func (s *Secret) Interweave(secret string) error {
	return errors.New("can't use interweave method on generic `secret-wrapper` class")
}
func (s *Secret) Unravel(path string) error {
	return errors.New("can't use unravel method on generic `secret-wrapper` class")
}
func (S *Secret) Data() *SecretData {
	return &S.data
}

func New(path string) SecretInterface {
	for _, media := range Media {
		if secret := media(path); secret != nil {
			return secret
		}
	}

	return nil
}
