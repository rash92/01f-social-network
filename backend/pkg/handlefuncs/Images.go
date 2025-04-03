package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
)

func SaveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// generate new uuid for image name
	uniqueId := uuid.New()
	// remove "- from imageName"
	filename := strings.Replace(uniqueId.String(), "-", "", -1)
	// extract image extension from original file filename
	fileExt := strings.Split(header.Filename, ".")[len(strings.Split(header.Filename, "."))-1]
	supported := isSupportedFileType(fileExt)

	if !supported {
		// rereturn "",error message to the user that this type of file is not supported
		return "", errors.New("this file type is not supported")
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

func isSupportedFileType(fileType string) bool {

	supportedTypes := map[string]bool{
		"jpeg": true,
		"jpg":  true,
		"png":  true,
		"gif":  true,
	}
	return supportedTypes[strings.ToLower(fileType)]
}

func ConvertBase64ToImage(base64String string) (*dbfuncs.File, error) {
	// Split the base64 string to isolate the MIME type and the actual data

	splitData := strings.Split(base64String, ",")
	if len(splitData) != 2 {
		return nil, fmt.Errorf("invalid base64 string")
	}

	mimeType := strings.Split(splitData[0], ";")[0]
	data := splitData[1]

	// Map the MIME type to a file extension
	mimeToExtension := map[string]string{
		"data:image/jpeg": ".jpg",
		"data:image/png":  ".png",
		"data:image/gif":  ".gif",
		// Add more mappings as needed
	}
	extension, ok := mimeToExtension[mimeType]
	if !ok {
		return nil, fmt.Errorf("unsupported file type: %s", mimeType)
	}

	// Decode the base64 string back to bytes
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	image := dbfuncs.File{
		Bytes:     decodedData,
		Extension: extension,
	}

	return &image, nil
}
