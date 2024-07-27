package image

import (
	"fmt"
	"github.com/zangarmarsh/inplainsight/core/steganography"
	"github.com/zangarmarsh/inplainsight/core/utility"
	"strings"
)

func init() {
	steganography.Media = append(
		steganography.Media,
		func(filePath string) steganography.SecretInterface {
			// ToDo: to ensure content integrity it might worth adding a manipulation (such as resizing) on a copy of the image
			// 			 since http.DetectContentType is not super reliable

			allegedContentType := utility.SniffMimeType(filePath)
			if strings.HasPrefix(allegedContentType, "image/") {
				fmt.Println("it is image")
				return NewImage(filePath)
			} else {
				return nil
			}
		},
	)
}
