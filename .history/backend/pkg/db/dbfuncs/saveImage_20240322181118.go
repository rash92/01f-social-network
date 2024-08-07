package dbfuncs

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

func SaveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// generate new uuid for image name

	// remove "- from imageName"
	filename := strings.Replace(uniqueId.String(), "-", "", -1)
	// extract image extension from original file filename
	fileExt := strings.Split(header.Filename, ".")[len(strings.Split(header.Filename, "."))-1]
	supported := isSupportedFileType(fileExt)

	if !supported {
		// rereturn "",error message to the user that this type of file is not supported
		return "", errors.New("this file type  is not supported")
	}

	// generate image from filename and extension
	image := fmt.Sprintf("%s.%s", filename, fileExt)
	// create a new file in the "uploads" folder
	dst, err := os.Create(fmt.Sprintf("pkg/db/images/%s", image))
	if err != nil {
		log.Println("unable to create file --> ", err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	return image, nil
}



func returnToStartingPoint(commands string) int {
	x, y := 0, 0
	direction := 0 // 0: North, 1: East, 2: South, 3: West

	for _, command := range commands {
		switch command {
		case 'F':
			switch direction {
			case 0:
				y++
			case 1:
				x++
			case 2:
				y--
			case 3:
				x--
			}
		}
			switch direction {
			case 0:
				y++
			case 1:
				x++
			case 2:
				y--
			case 3:
				x--
			}
		case 'L':
			direction = (direction + 3) % 4
		case 'R':
			direction = (direction + 1) % 4
		}
	}

	return abs(x) + abs(y)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}