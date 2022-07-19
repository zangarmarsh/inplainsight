package steganography

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
)

type steganography struct {
	header header
	image  image.Image
}

type header struct {
	version      [5]byte
	magicNumber  [13]byte
	compression  uint8
	endOfChar    uint8
	endOfMessage uint8
}

func (s *steganography) SetHeader( message string, maximumCompression, endOfChar, endOfMessage uint8 ) error {
	if maximumCompression < 1 || maximumCompression > 8 {
		return errors.New( "compression level must be between 1 and 8" )
	}

	h := header{ }

	h.endOfMessage = endOfMessage
	h.endOfChar    = endOfChar

	imageSize := s.image.Bounds().Size()

	compression := estimateCompressionLevel( imageSize.X * imageSize.Y, []byte(message) )
	if maximumCompression > compression {
		return errors.New( "the minimum compression level for this message and image is " )
	}

	s.header = h
	return nil
}

func estimateCompressionLevel( amountOfPixels int, message []byte ) uint8 {
	// 3*p*h*w/(v+3*p*l+3*p)
	// channels := 3 // R G B
	// cumulableValue := uint64( amountOfPixels * channels )

	messageValue := uint64(0)

	for v, _ := range message {
		messageValue += uint64( v )
	}

	// intermediateRequiredSpace := messageValue / cumulableValue
  // intermediateRequiredSpace

	return 3
}

func Reveal( in string, loss uint8 ) (string, error) {
	var secretMessage []byte

	img, err := getImageContent( in )
	if err != nil {
		return "", err
	}

	width, height := img.Bounds().Size().X, img.Bounds().Size().Y
	bitMask := ^uint8(0) << loss

	var lastChar uint8
	var bufChar uint8

	for y := 0; y < height; y++  {
		for x := 0; x < width; x++ {
			r,g,b,_ := img.At( x,y ).RGBA( )

			redOffset   := ^bitMask & uint8( r )
			greenOffset := ^bitMask & uint8( g )
			blueOffset  := ^bitMask & uint8( b )

			bufChar  += redOffset + greenOffset + blueOffset

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

	return string( secretMessage ), nil
}

func Conceal( in, out, secretMessage string, loss uint8 ) error {
	fmt.Printf( "Input file: %s\nOutput file: %s\nSecret message: %s\n", in, out, secretMessage )

	img, err := getImageContent( in )
	if err != nil {
		return err
	}

	outImage := image.NewRGBA( img.Bounds() )

	x, y, err := conceal( outImage, secretMessage, img, loss )

	if err != nil {
		return err
	}

	// @ToDo Make it stronger
	for i := 1; i <= 2; i++ {
		r,g,b,a := img.At( x+i, y ).RGBA( )

		// Adds a final marker
		outImage.Set( x+i, y, color.RGBA {
			G: uint8(g) & ( ^uint8(0) << loss ),
			R: uint8(r) & ( ^uint8(0) << loss ),
			B: uint8(b) & ( ^uint8(0) << loss ),
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

func conceal( outImage *image.RGBA, secretMessage string, img image.Image, loss uint8 ) ( x, y int, err error ) {
	err = nil
	width, height := outImage.Bounds().Size().X, outImage.Bounds().Size().Y

	maxValue := uint8(1 << loss) - 1
	secretChar := secretMessage[0]
	secretChars := secretMessage[1:]

	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			compensations := make( []uint8, 3 )
			r, g, b, a := img.At(x, y).RGBA( )

			if secretChar != 0 {
				for compensation := uint8(0); compensation < 3; compensation++ {
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
					R: uint8(r) &bitmask | compensations[0],
					G: uint8(g) &bitmask | compensations[1],
					B: uint8(b) &bitmask | compensations[2],
					A: uint8(a),
				})
		}
	}

	return
}

func getImageContent( in string ) ( image.Image, error ) {
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
	return img, nil
}