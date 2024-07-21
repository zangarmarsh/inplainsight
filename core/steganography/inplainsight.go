package steganography

import (
	_ "image/jpeg"
	_ "image/png"
)

const Version uint8 = '\x01'
const MagicNumber uint8 = '\x78'
const EndOfMessage uint8 = '\x00'
