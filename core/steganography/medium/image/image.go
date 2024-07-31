package image

import (
	"errors"
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"unsafe"
)

const bitsPerChannel int = 3
const channelsPerPixel int = 3

type Image struct {
	steganography.Host
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
		uint64(i.Header.Bits())) /
		uint64(unsafe.Sizeof(uint32(0))) * 8
}

// Len Counts the UTF-8 characters currently interwoven
func (i *Image) Len() uint64 {
	return uint64(len(i.Data().Encrypted))
}

func (i *Image) Interweave(secret string) error {
	if len(secret) == 0 {
		return errors.New("Cannot interweave empty secret")
	}

	// Todo: dynamic determination of the optimal "bits per channel" amount by the overall image size
	bitsPerChannel := uint8(bitsPerChannel)
	bitmask := ^uint8(math.Pow(2, float64(bitsPerChannel)) - 1)

	{
		headerPtr, err := steganography.NewHeader(bitsPerChannel, steganography.EndOfMessage)
		if err != nil {
			return err
		}
		i.Header = headerPtr
		i.Header.Compression = bitsPerChannel
	}

	// Check if resource is available and properly set up
	if i.resource == nil {
		return errors.New("interweave called before a proper initialization, retry calling `NewImage` method first")
	}

	// You'll need some sewing thread to interweave stuff into the picture :D
	yarn, err := i.CraftYarn(secret)
	if err != nil {
		return err
	}

	// Check if the image has enough space to store the message
	if uint64(len(secret)) > i.Cap()-i.Len() {
		return errors.New("secret too long")
	}

	width, height := (*i.resource).Bounds().Size().X, (*i.resource).Bounds().Size().Y
	output := image.NewNRGBA(image.Rect(0, 0, width, height))

	{
		bitsChan := make(chan uint8)

		go steganography.CutYarnChunks(bitsChan, yarn, int(bitsPerChannel))

		cloneExistingPixel := false
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, _ := (*i.resource).At(x, y).RGBA()

				// Values must be right-shifted before casting because otherwise it would just keep the last - potentially random - eight bits,
				// returning way less precise colors.
				red := uint8(r >> 8)
				green := uint8(g >> 8)
				blue := uint8(b >> 8)

				for _, c := range []*uint8{&red, &green, &blue} {
					if !cloneExistingPixel {
						bits, ok := <-bitsChan

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
	}

	outFile, err := os.Create((*i).Path)
	if err != nil {
		return err
	}

	defer outFile.Close()

	err = png.Encode(outFile, output)
	if err != nil {
		return err
	}

	i.Data().Encrypted = secret

	return nil
}

func (i *Image) Unravel(path string) error {
	// Open up file handle if did not already
	if i.Path == "" || i.resource == nil {
		err := i.setImage(path)
		if err != nil {
			return err
		}
	}

	i.Header = &steganography.Header{}

	// Let's try to decrypt the header
	{
		unraveled := make([]rune, 0)
		width, height := (*i.resource).Bounds().Size().Y, (*i.resource).Bounds().Size().Y

		// This could escalate quickly, evaluate if worths using a huge number instead
		writtenBits := int64(0)

		(func() {
			iterationCounter := 0
			initialOffset := int(unsafe.Sizeof('\x00')) * 8
			offset := initialOffset
			iterationsForOneByte := int(math.Ceil(float64(unsafe.Sizeof(uint8(0))*8) / float64(bitsPerChannel)))

			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					r, g, b, _ := (*i.resource).At(x, y).RGBA()

					for _, c := range []uint32{r, g, b} {
						bitmask := rune(uint8(math.Pow(2, float64(bitsPerChannel)) - 1))

						// It should automatically truncate the real number part but let's ensure that
						characterIndex := int(math.Floor(float64(writtenBits / (int64(unsafe.Sizeof('\x00')) * 8))))

						offset -= bitsPerChannel

						// todo: what is this shit doing?
						if characterIndex >= len(unraveled) {
							unraveled = append(unraveled, '\x00')
						}

						if adjustment := 8 - (iterationCounter%iterationsForOneByte)*bitsPerChannel - bitsPerChannel; adjustment < 0 {
							adjustment := int(math.Abs(float64(adjustment)))

							bitmask >>= adjustment
							offset += adjustment
						}

						if offset < 0 {
							offset = 0
						}

						unraveled[characterIndex] |= bitmask & int32(uint8(c)) << offset

						iterationCounter++
						writtenBits += int64(math.Log2(float64(bitmask + 1)))

						if offset == 0 {
							if i.Header.MagicNumber != 0 && unraveled[characterIndex] == rune(steganography.EndOfMessage) {
								unraveled = unraveled[1:]
								unraveled = unraveled[:len(unraveled)-1]

								return
							}

							if i.Header.MagicNumber == 0 {
								magicNumberMatched := i.Header.Set(unraveled[0])

								if !magicNumberMatched {
									unraveled = nil
									return
								}
							}

							offset = initialOffset
						}
					}
				}
			}

			return
		})()

		i.Data().Encrypted = string(unraveled)
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

	i.Path = path

	i.resource = &img
	return nil
}
