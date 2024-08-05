package image_test

import (
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"golang.org/x/exp/rand"
	"os/exec"
	"testing"

	_ "embed"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const blankSampleFile = "test/samples/test.jpg"

// Parameter `size` must have this structure: [width]x[height]
func generateBlankImage(size string) error {
	if err := exec.Command("rm", "rm", blankSampleFile).Run(); err != nil {
		return err
	}

	command := exec.Command(
		"convert",
		"convert",
		"-size",
		size,
		"xc:white",
		blankSampleFile,
	)

	return command.Run()
}

func generateText(size int) string {
	var text string

	for c := 0; c < size; c++ {
		char := rand.Int31()
		if char == 0 {
			char++
		}
		text += string(char)
	}

	return text
}

func TestSteganography(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Steganography Suite")
}

var _ = Describe("Concealing/Revealing", func() {
	Context("When using a non existant image", func() {
		secret := steganography.New("/invalid/path/file")
		It("Should stops gracefully", func() {
			Expect(secret).To(BeNil())
		})
	})

	Context("When a valid png image is concealed", func() {
		generateBlankImage("100x100")
		text := "私は inplainsight です!!"

		It("Stops since there's no text to conceal", func() {
			secret := steganography.New(blankSampleFile)
			Expect(secret).NotTo(BeNil())

			err := secret.Interweave("")
			Expect(err).ShouldNot(BeNil())
		})
		//

		It("Conceals text into a test sample", func() {
			s := steganography.New(blankSampleFile)
			Expect(s).ShouldNot(BeNil())

			err := s.Interweave(text)
			Expect(err).To(BeNil())
		})

		It("Reveals previously concealed text from a test sample", func() {
			s := steganography.New(blankSampleFile)
			Expect(s).ShouldNot(BeNil())

			Expect(*s.Data()).To(BeEquivalentTo(text))
		})
	})
})
