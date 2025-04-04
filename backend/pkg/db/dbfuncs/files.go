package dbfuncs

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func SaveImage(imageFile File) (string, error) {
	imageId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	fileName := imageId.String() + imageFile.Extension

	imagePath := filepath.Join(imageDirectory, fileName)
	err = os.WriteFile(imagePath, imageFile.Bytes, 0644)
	if err != nil {
		return "", err
	}
	return fileName, nil
}
