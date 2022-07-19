package steganography_test

import (
	"testing"
	"webshapes.it/steganography"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSteganography(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Steganography Suite")
}

var _ = Describe("Concealing/Revealing", func() {
	loss := uint8(3)
	message := "Dammi i bitcoin!!"
	in := "in/test.png"
	out := "out/test.png"

	Context("When a valid png image is concealed losing the least three significant bits for each RGB(A) matrix", func( ) {
		err := steganography.Conceal(in, out, message, loss )

		It("Conceils successfully", func() {
			Expect( err ).Should( BeNil() )
		})

		// exec.Command("xdg-open", "./out/test.png").Run()
	})

	Context("When a png image is revealed specifying to use the least three significant bits", func() {
		revealed, err := steganography.Reveal(out, loss)

		It("Reveals the message", func() {
			Expect(revealed).Should( Equal(message) )
		})

		It("Should not have any error", func() {
			Expect(err).Should( BeNil() )
		})
	})
})