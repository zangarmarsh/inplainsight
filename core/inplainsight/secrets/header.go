package secrets

type Header struct {
	mn            MagicNumber
	dedicatedHost bool
}

func NewHeader(bitmask byte) *Header {
	firstBitMask := byte(1 << 7)

	return &Header{
		dedicatedHost: firstBitMask&bitmask > 0,
		mn:            MagicNumber(^firstBitMask & bitmask),
	}
}
