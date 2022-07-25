package steganography_test

import (
	"testing"
	"webshapes.it/steganography"

	_ "embed"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:embed samples/texts/short.txt
var shortText []byte // can be handled with a compression of one
//go:embed samples/texts/normal.txt
var normalText []byte // 3 bits of compression required
//go:embed samples/texts/big.txt
var bigText []byte // can't be compressed in the input image

func TestSteganography(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Steganography Suite")
}

var _ = Describe("Concealing/Revealing", func() {
	loss := uint8(3)
	in := "samples/in/test.png"
	out := "samples/out/test.png"

	BeforeEach(func() {
		s := steganography.Steganography{}
		Expect(s).ToNot(BeNil())
	})

	Context("When a valid png image is concealed", func( ) {
		By("losing the least three significant bits for each RGB(A) matrix")

		It("Conceals successfully", func() {
			s := new(steganography.Steganography)

			err := s.Conceal( in, out, string(normalText), loss )
			Expect( err ).Should( BeNil() )
		})

		It("Stops if the required compression is way higher than the indicated one", func() {
			s := new(steganography.Steganography)

			err := s.Conceal(in, out, string(bigText), 3)
			Expect( err ).Should( Not( BeNil() ))
		})
	})

	Context("When a png image is revealed", func() {
		By( "specifying to use the least three significant bits" )

		It("Reveals the message", func() {
			s := new( steganography.Steganography )
			err := s.Conceal( in, out, string(normalText), loss )

			revealed, err := s.Reveal( out )
			Expect( revealed ).Should( Equal(string(normalText)) )
			Expect( err ).Should( BeNil() )
		})
	})
})
Describe( "When a valid unconcealed png is revealed", Pending, )
