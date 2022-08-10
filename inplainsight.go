package inplainsight

import (
	"errors"
	"fmt"
	"github.com/zangarmarsh/inplainsight/cryptography"
	"github.com/zangarmarsh/inplainsight/header"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"math"
	"os"
	"reflect"
)

const version uint8 									 = '\x01'
const magicNumber uint8 							 = '\x78'
const inweavedHeaderBitsPerChannel int = 2
const inweaveableChannelsPerPixel int  = 3


type Steganography struct {
	header header.Header
	image  *image.Image
}

func (s *Steganography) AmountOfSkippablePixels() int {
	return int(math.Ceil(float64(s.header.Size()) / float64(inweavedHeaderBitsPerChannel) / float64(3)))
}

func (s *Steganography) SetHeader(message string, maximumCompression, endOfChar, endOfMessage uint8) error {
	if maximumCompression < 1 || maximumCompression > 8 {
		return errors.New("compression level must be between 1 and 8")
	}

	s.header.MagicNumber = magicNumber
	s.header.EndOfCharDelimiter = endOfMessage
	s.header.EndOfMessageDelimiter = endOfChar

	imageSize := (*s.image).Bounds().Size()

	compression, err := estimateCompressionLevel(
		uint64(imageSize.X * imageSize.Y) -
		uint64(math.Ceil(float64(s.header.Size()) / float64(inweaveableChannelsPerPixel) /
		float64(inweavedHeaderBitsPerChannel))),

		[]byte(message),
	)
	if err != nil {
		return errors.New("this message cannot be hid, try again with a bigger picture")
	}
	if compression > maximumCompression {
		return errors.New(fmt.Sprintf("the minimum compression level for this message and image is %d, %d found", maximumCompression, compression))
	}

	log.Printf("It's been estimated a compression level of %d/8\n", compression)

	s.header.Version = version
	s.header.Compression = compression

	return nil
}

func estimateCompressionLevel(amountOfPixels uint64, message []byte) (uint8, error) {
	var messageValue uint64 = 0

	for _, v := range message {
		messageValue += uint64(v)
	}

	messageValue += uint64(len(message)) + 2

	return uint8(
		math.Max(
			1,
			math.Ceil(
				1 / ( float64( amountOfPixels * uint64(inweaveableChannelsPerPixel) ) / float64(messageValue) ),
			),
		),
	),
	nil
}

func (s *Steganography) Reveal(in, password string) (string, error) {
	var secretMessage []byte
	skipPixels := s.AmountOfSkippablePixels()

	img, err := getImageContent(in)
	if err != nil {
		return "", err
	}

	err = s.extractHeader(img)
	if err != nil {
		return "", err
	}

	width, height := (*img).Bounds().Size().X, (*img).Bounds().Size().Y
	bitMask := ^uint8(0) << s.header.Compression

	var lastChar uint8
	var bufChar uint8

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (y+1)*(x+1) <= skipPixels {
				continue
			}

			r, g, b, _ := (*img).At(x, y).RGBA()

			redOffset   := ^bitMask & uint8(r)
			greenOffset := ^bitMask & uint8(g)
			blueOffset  := ^bitMask & uint8(b)

			bufChar += redOffset + greenOffset + blueOffset

			if redOffset == 0 || greenOffset == 0 || blueOffset == 0 {
				if bufChar == 0 && lastChar == 0 {
					return string(secretMessage), nil
				} else if bufChar != 0 {
					secretMessage = append(secretMessage, bufChar)
					lastChar = bufChar
					bufChar = 0
				}
			}
		}
	}
	if len(password) != 0 {
		contentEncriptionKey, headerEncriptionKey := cryptography.DeriveEncryptionKeysFromPassword(password)
		_ = headerEncriptionKey // ToDo: remove it when the header will be encrypted

		decryptedMessage, err := cryptography.Decrypt(secretMessage, contentEncriptionKey)
		if err != nil {
			return "", err
		}

		secretMessage = decryptedMessage
	}

	return string(secretMessage), nil
}

func (s *Steganography) Conceal(in, out string, secretMessage []byte, password string, maximumCompression uint8) error {
	if len(secretMessage) == 0 {
		return errors.New("The provided message is empty" )
	}

	//fmt.Printf("Input file: %s\nOutput file: %s\nSecret message: %s...\n", in, out, secretMessage[:25])

	img, err := getImageContent(in)
	if err != nil {
		return err
	}

	s.image = img

	imgSize := (*img).Bounds()
	outImage := image.NewRGBA(imgSize)
	delimiter := uint8(0)

	if len(password) != 0 {
		contentEncriptionKey, headerEncriptionKey := cryptography.DeriveEncryptionKeysFromPassword(password)
		_ = headerEncriptionKey // ToDo: remove it when the header will be encrypted

		secretMessage, err = cryptography.Encrypt(secretMessage, contentEncriptionKey)
		if err != nil {
			return err
		}
	}

	if err = s.SetHeader(string(secretMessage), maximumCompression, delimiter, delimiter); err != nil {
		return err
	}

	x, y, err := conceal(outImage, secretMessage, img, s.header.Compression, s.AmountOfSkippablePixels())
	if err != nil {
		return err
	}

	s.interweaveHeader(outImage)

	// @ToDo Make it stronger
	for i := 1; i <= 2; i++ {
		r, g, b, a := (*img).At(x+i, y).RGBA()

		// Adds a final marker
		outImage.Set(x+i, y, color.RGBA{
			G: uint8(g) & (^uint8(0) << s.header.Compression),
			R: uint8(r) & (^uint8(0) << s.header.Compression),
			B: uint8(b) & (^uint8(0) << s.header.Compression),
			A: uint8(a),
		})
	}

	outFile, err := os.Create(out)
	if err != nil {
		return err
	}

	defer outFile.Close()

	err = png.Encode(outFile, outImage)
	if err != nil {
		return err
	}

	return nil
}

func (s *Steganography) interweaveHeader(outImage *image.RGBA) {
	size := (*outImage).Bounds().Size()
	h := []uint8{
		s.header.MagicNumber,
		s.header.Version,
		s.header.Compression,
		s.header.EndOfCharDelimiter,
		s.header.EndOfMessageDelimiter,
	}

	additionBitmask := uint8(math.Pow(float64(2), float64(inweavedHeaderBitsPerChannel)) - 1)
	shiftableBitmask := additionBitmask << (8 - inweavedHeaderBitsPerChannel)
	blocks := int(math.Ceil(float64(s.header.Size()) / float64(inweavedHeaderBitsPerChannel) / float64(inweaveableChannelsPerPixel)))
	var fieldIndex int

	for i := 0; i < blocks; i++ {
		x, y := i%size.Y, i/size.Y

		additions := make([]uint8, inweaveableChannelsPerPixel)
		pixel := outImage.At(x, y)

		colors := make([]uint32, inweaveableChannelsPerPixel)
		colors[0], colors[1], colors[2], _ = pixel.RGBA()

		{
			for splitting := 0; splitting < cap(additions); splitting++ {
				if int( math.Floor(float64(splitting+(i*cap(additions)))/4) ) >= s.header.Size() / 8 {
					break
				}

				fieldIndex = int(math.Floor(float64(splitting+(i*cap(additions))) / 4))
				//fmt.Printf( "[based on value %08b]\n", h[fieldIndex])

				offset := (splitting + (i * cap(additions))) % 4 * 2
				shiftedBitmask := shiftableBitmask >> offset
				additions[splitting] = (shiftedBitmask & h[fieldIndex]) >> (6 - offset)
				//fmt.Printf("addition %08b\n", additions[splitting])
			}
		}

		(*outImage).Set(x, y, color.RGBA{
			R: (uint8(colors[0]) & additionBitmask) | additions[0],
			G: (uint8(colors[1]) & additionBitmask) | additions[1],
			B: (uint8(colors[2]) & additionBitmask) | additions[2],
			A: uint8(255),
		})
	}

}

func (s *Steganography) extractHeader(img *image.Image) error {
	size := (*img).Bounds().Size()
	headerSize := s.header.Size()
	colors := make([]uint32, inweaveableChannelsPerPixel)

	fields := make([]byte, reflect.TypeOf(s.header).NumField())
	bitmask := uint8(math.Pow(2, float64(inweavedHeaderBitsPerChannel)) - 1)
	pixelsForHeader := int(math.Ceil(float64(headerSize) / float64(inweavedHeaderBitsPerChannel) / float64(cap(colors))))

	var currentPixel = 0

	for index := 0; index < pixelsForHeader; index++ {
		colors[0], colors[1], colors[2], _ = (*img).At(index%size.Y, index/size.Y).RGBA()

		for channelIndex := 0; channelIndex < cap(colors) && (currentPixel*inweavedHeaderBitsPerChannel) < headerSize; channelIndex++ {
			fieldIndex := int(math.Floor(float64(currentPixel * inweavedHeaderBitsPerChannel) / 8))
			currentPixel++

			info := bitmask & uint8(colors[channelIndex])
			amountToShift := 6 - ((channelIndex + (cap(colors) * index)) % 4 * inweavedHeaderBitsPerChannel)

			fields[fieldIndex] += info << amountToShift
		}
	}

	s.header = header.Header{
		MagicNumber: fields[0],
		Version: fields[1],
		Compression: fields[2],
		EndOfCharDelimiter: fields[3],
		EndOfMessageDelimiter: fields[4],
	}

	if s.header.MagicNumber != magicNumber {
		return errors.New( "the given image either not concealed through steganography or needs a password" )
	}

	return nil
}

func conceal(outImage *image.RGBA, secretMessage []byte, img *image.Image, loss uint8, skipPixels int) (x, y int, err error) {
	err = nil
	width, height := outImage.Bounds().Size().X, outImage.Bounds().Size().Y

	maxValue := uint8(1<<loss) - 1
	secretChar := secretMessage[0]
	secretChars := secretMessage[1:]

	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			if (x+1)*(y+1) <= skipPixels {
				bitmask := ^uint8(math.Pow(2, float64(inweavedHeaderBitsPerChannel)) - 1)
				r,g,b,a := (*img).At(x,y).RGBA()
				outImage.Set(x,y, color.RGBA{
					R: uint8(r) & bitmask,
					G: uint8(g) & bitmask,
					B: uint8(b) & bitmask,
					A: uint8(a),
				})
				continue
			}

			compensations := make([]uint8, inweaveableChannelsPerPixel)
			r, g, b, a := (*img).At(x, y).RGBA()

			if secretChar != 0 {
				for compensation := 0; compensation < inweaveableChannelsPerPixel; compensation++ {
					var value uint8

					if secretChar < maxValue {
						value = secretChar
					} else {
						value = maxValue
					}

					secretChar -= value
					compensations[compensation] = value
				}
			} else if len(secretChars) > 0 {
				secretChar = secretChars[0]
				secretChars = secretChars[1:]
			}

			bitmask := ^uint8(0) << loss
			outImage.Set(
				x,
				y,
				color.RGBA{
					R: uint8(r)&bitmask | compensations[0],
					G: uint8(g)&bitmask | compensations[1],
					B: uint8(b)&bitmask | compensations[2],
					A: uint8(a),
				})
		}
	}

	return
}

func getImageContent(in string) (*image.Image, error) {
	fileContent, err := os.Open(in)
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	_, err = fileContent.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(fileContent)
	if err != nil {
		return nil, err
	}
	return &img, nil
}
