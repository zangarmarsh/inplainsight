package steganography

type MediumRegistrator func(filePath string) SecretInterface

var Media []MediumRegistrator
