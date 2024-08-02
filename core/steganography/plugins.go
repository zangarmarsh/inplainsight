package steganography

type MediumRegistrator func(filePath string) SecretHostInterface

var Media []MediumRegistrator
