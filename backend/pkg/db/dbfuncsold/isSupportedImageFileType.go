package dbfuncs

import "strings"

func isSupportedFileType(fileType string) bool {

	supportedTypes := map[string]bool{
		"jpeg": true,
		"jpg":  true,
		"png":  true,
		"gif":  true,
	}
	return supportedTypes[strings.ToLower(fileType)]
}
