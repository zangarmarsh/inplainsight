package image_test

import (
	"github.com/zangarmarsh/inplainsight/core/steganography/medium/image"
	"os/exec"
	"testing"

	_ "embed"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const blankSampleFile = "test/samples/blank.png"

// //go:embed ../../test/texts/short.txt
// var shortText []byte // can be handled with a compression of one
// //go:embed samples/texts/normal.txt
// var normalText []byte // 3 bits of compression required
// //go:embed samples/texts/big.txt
// var bigText []byte // can't be compressed in the input image

// Parameter `size` must have this structure: [width]x[height]
func generateBlankImage(size string) error {
	if err := exec.Command("rm", "rm", blankSampleFile).Run(); err != nil {
		command := exec.Command(
			"convert",
			"convert",
			"-size",
			size,
			"xc:white",
			blankSampleFile,
		)

		return command.Run()
	} else {
		return err
	}
}

func TestSteganography(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Steganography Suite")
}

var _ = Describe("Concealing/Revealing", func() {
	Context("When using a non existant image", func() {
		secret := image.NewImage("/invalid/path/file")
		It("Should stops gracefully", func() {
			Expect(secret).To(BeNil())
		})
	})

	Context("When a valid png image is concealed", func() {
		generateBlankImage("20x20")
		text := "私は inplainsight です!!"

		It("Stops since there's no text to conceal", func() {
			secret := image.NewImage(blankSampleFile)
			err := secret.Interweave("")
			Expect(err).ShouldNot(BeNil())
		})
		//

		It("Conceals text into a test sample", func() {
			s := image.NewImage(blankSampleFile)
			Expect(s).ShouldNot(BeNil())

			err := s.Interweave(text)
			Expect(err).To(BeNil())
		})

		It("Reveals previously concealed text from a test sample", func() {
			s := image.NewImage(blankSampleFile)
			Expect(s).ShouldNot(BeNil())

			Expect(s.Data.Decrypted).To(BeEquivalentTo(text))
		})
	})
})
