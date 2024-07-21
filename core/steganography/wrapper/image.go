package wrapper

import (
	"errors"
	"fmt"
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"github.com/zangarmarsh/inplainsight/core/steganography/header"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"unsafe"
)

const bitsPerChannel int = 2
const channelsPerPixel int = 3

type Image struct {
	SecretWrapper
	resource *image.Image
}

func NewImage(imagePath string) *Image {
	img := Image{}
	err := img.setImage(imagePath)

	if err != nil {
		log.Println(err)
		return nil
	}

	if img.Unravel(imagePath) == nil {
		return &img
	} else {
		return nil
	}
}

// Get the count of how many UTF-8 characters can be interwoven within an image
func (i *Image) Cap() uint64 {
	bounds := (*i.resource).Bounds()

	return (uint64(bounds.Dx())*uint64(bounds.Dy())*
		uint64(bitsPerChannel)*
		uint64(channelsPerPixel) -
		uint64(i.header.Bits())) /
		uint64(unsafe.Sizeof(uint32(0))) * 8
}

// Counts the UTF-8 characters currently interwoven
func (i *Image) Len() uint64 {
	return uint64(len(i.data.encrypted))
}

func (i *Image) Interweave(secret string) error {
	if len(secret) == 0 {
		return errors.New("Cannot interweave empty secret")
	}

	// Todo: dinamically determination of the optimal "bits per channel" amount by the overall image size
	bitsPerChannel := uint8(bitsPerChannel)
	bitmask := ^uint8(math.Pow(float64(bitsPerChannel), 2) - 1)

	{
		headerPtr, err := header.NewHeader(bitsPerChannel, steganography.EndOfMessage)
		if err != nil {
			return err
		}
		i.header = headerPtr
		i.header.Compression = bitsPerChannel
	}

	// Check if resource is available and properly set up
	if i.resource == nil {
		return errors.New("interweave called before a proper initialization, retry calling `NewImage` method first")
	}

	// First of all, check if the image has enough space to store the message
	if uint64(len(secret)) > i.Cap()-i.Len() {
		return errors.New("secret too long")
	}

	// You'll need some sewing thread to interweave stuff into the picture :D
	yarn, err := i.craftYarn(secret)
	if err != nil {
		return err
	}

	fmt.Println("gonna conceal", yarn)

	width, height := (*i.resource).Bounds().Size().X, (*i.resource).Bounds().Size().Y
	fmt.Println("width", width, "height", height)
	output := image.NewRGBA(image.Rect(0, 0, width, height))

	{
		bits := make(chan uint8)

		go cutYarnChunks(bits, yarn, int(bitsPerChannel))

		(func() {
			cloneExistingPixel := false
			for y := 0; y <= height; y++ {
				for x := 0; x <= width; x++ {
					r, g, b, _ := (*i.resource).At(x, y).RGBA()
					red := uint8(r)
					green := uint8(g)
					blue := uint8(b)

					for _, c := range []*uint8{&red, &green, &blue} {
						if !cloneExistingPixel {
							bits, ok := <-bits
							if ok {
								*c = (*c & bitmask) | bits
							} else {
								cloneExistingPixel = true
							}
						}
					}

					output.Set(x, y, color.RGBA{
						R: red,
						G: green,
						B: blue,
						A: ^uint8(0),
					})
				}
			}
		})()
	}

	outFile, err := os.Create((*i).path)
	if err != nil {
		return err
	}

	defer outFile.Close()

	err = png.Encode(outFile, output)
	if err != nil {
		return err
	}

	return nil
}

func (i *Image) Unravel(path string) error {
	// Open up file handle if did not already
	if i.path == "" || i.resource == nil {
		err := i.setImage(path)
		if err != nil {
			return err
		}
	}

	// Let's try to decrypt the header
	{
		unraveled := make([]rune, 0)
		width, height := (*i.resource).Bounds().Size().Y, (*i.resource).Bounds().Size().Y

		bitmask := uint8(math.Pow(2, float64(bitsPerChannel)) - 1)
		iterationsPerRune := int(unsafe.Sizeof(rune(0))) * 8 / bitsPerChannel

		(func() {
			iterationCounter := 0
			for y := 0; y <= height; y++ {
				for x := 0; x <= width; x++ {
					r, g, b, _ := (*i.resource).At(x, y).RGBA()

					for _, c := range []uint32{r, g, b} {
						// It should automatically truncate the real number part but let's ensure that
						characterIndex := int(math.Floor(float64(iterationCounter) / float64(iterationsPerRune)))
						groupOfBitsIndex := iterationCounter % iterationsPerRune

						if characterIndex >= len(unraveled) {
							unraveled = append(unraveled, '\x00')
						}

						unraveled[characterIndex] |= int32(uint8(c)) & rune(bitmask) << (int(unsafe.Sizeof('\x00'))*8 - (groupOfBitsIndex+1)*bitsPerChannel)

						iterationCounter++

						if groupOfBitsIndex == iterationsPerRune-1 &&
							unraveled[characterIndex] == rune(steganography.EndOfMessage) {
							return
						}
					}
				}
			}
		})()

		i.data.decrypted = string(unraveled)
		fmt.Printf("UNRAVELED %+v\n%+v\n", string(unraveled), unraveled)
	}

	return nil
}

func (i *Image) setImage(path string) error {
	fileContent, err := os.Open(path)
	if err != nil {
		return err
	}

	defer fileContent.Close()

	_, err = fileContent.Seek(0, 0)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(fileContent)
	if err != nil {
		return err
	}

	i.path = path

	i.resource = &img
	return nil
}
