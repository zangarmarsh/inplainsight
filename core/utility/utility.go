package utility

import (
	"net/http"
	"os"
)

func SniffMimeType(filePath string) string {
	// ToDo: to ensure content integrity it might worth adding a manipulation (such as resizing) on a copy of the image
	// 			 since http.DetectContentType is not super reliable
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
