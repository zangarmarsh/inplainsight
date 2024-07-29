package image

import (
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"github.com/zangarmarsh/inplainsight/core/utility"
	"log"
	"os"
	"strings"
)

func init() {
	steganography.Media = append(
		steganography.Media,
		func(filePath string) steganography.SecretInterface {
			if _, err := os.Stat(filePath); err != nil {
				log.Printf("File %v does not exist", filePath)
				return nil
			}

			if allegedContentType := utility.SniffMimeType(filePath); strings.HasPrefix(allegedContentType, "image/") {
				return NewImage(filePath)
			} else {
				return nil
			}
		},
	)
}
