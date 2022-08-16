package steganography_test

import (
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"testing"

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
	maximumCompression := uint8(3)
	in       := "samples/in/test.png"
	out      := "samples/out/test.png"
	password := []byte("password")

	BeforeEach(func() {
		s := steganography.Steganography{}
		Expect(s).ToNot(BeNil())
	})

	Context("When a valid png image is concealed", func( ) {
		By("losing the least three significant bits for each RGB(A) matrix")

		It("Stops if there's no text to conceal", func() {
			s := new(steganography.Steganography)

			err := s.Conceal(in, out, nil, password, maximumCompression)
			Expect(err).Should( Not(BeNil()) )
		})

		It("Conceals successfully", func() {
			s := new(steganography.Steganography)

			err := s.Conceal(in, out, shortText, password, maximumCompression)
			Expect(err).Should( BeNil() )
		})

		It("Stops if the required compression is way higher than the indicated one", func() {
			s := new(steganography.Steganography)

			err := s.Conceal(in, out, bigText, password, 3)
			Expect(err).Should( Not( BeNil() ))
		})
	})

	Context("When a png image is revealed", func() {
		By( "specifying to use the least three significant bits" )

		It("Reveals the message", func() {
			s := new(steganography.Steganography)
			err := s.Conceal(in, out, shortText, nil, maximumCompression)

			revealed, err := s.Reveal(out, nil )
			Expect(revealed).Should( Equal(string(shortText)) )
			Expect(err).Should( BeNil() )
		})
	})

	Context("When a png image is revealed using a password", func() {
		By( "specifying to use the least three significant bits" )

		It("Reveals the message", func() {
			s := new(steganography.Steganography)
			err := s.Conceal(in, out, shortText, password, maximumCompression)

			revealed, err := s.Reveal(out, password)
			Expect(revealed).Should( Equal(string(shortText)) )
			Expect(err).Should( BeNil() )
		})
	})
})
