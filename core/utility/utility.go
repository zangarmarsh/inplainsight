package utility

import (
	"net/http"
	"os"
)

func SniffMimeType(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}

	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return ""
	}

	mimeType := http.DetectContentType(buffer)
	return mimeType
}
