package steganography

type MediumRegistrator func(filePath string) HostInterface

var Media []MediumRegistrator
