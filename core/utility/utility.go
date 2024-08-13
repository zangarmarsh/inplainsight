package utility

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
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

// SuggestFSPath tries to provide autocomplete suggestions based on a search query.
// If the specified file or folder already exists, it returns the same string. If there is a suitable suggestion, it
// returns the suggested string. Otherwise, if the query does not match any valid path and no suggestible file or folder
// is found, SuggestFSPath returns an empty string.
func SuggestFSPath(query string) (suggestion string) {
	if strings.ToLower(runtime.GOOS) == "linux" || strings.ToLower(runtime.GOOS) == "darwin" {
		if query[0] == '~' {
			homePath, _ := os.UserHomeDir()

			if homePath != "" {
				query = homePath + query[1:]

				defer func() {
					if suggestion != "" {
						suggestion = "~" + suggestion[len(homePath):]
					}
				}()
			}
		}
	}

	// If the query matches an existing file just get back
	if stat, err := os.Stat(query); err == nil {
		if stat != nil && (stat.Mode().IsRegular() || stat.Mode().IsDir()) {
			suggestion = query
			return
		}
	}

	var folder, suffix string
	if index := strings.LastIndex(query, string(os.PathSeparator)); index != -1 {
		folder = query[:index]

		if index != len(query)-1 {
			// Suffix does not contain the leftmost os.PathSeparator
			suffix = query[index+1:]
		}
	}

	if len(suffix) >= 1 {
		entries, err := os.ReadDir(folder)
		if err != nil {
			return
		}

		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), suffix) {
				suggestion = folder + string(os.PathSeparator) + entry.Name()

				if entry.IsDir() {
					suggestion += string(os.PathSeparator)
				}

				return
			}
		}
	}

	// The given path is wrong, ring the bell
	fmt.Printf("\a")

	return
}
