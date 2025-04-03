package dbfuncs

// func ConvertBase64ToImage(base64String string, directoryPath string) (string, error) {
// 	// Split the base64 string to isolate the MIME type and the actual data

// 	splitData := strings.Split(base64String, ",")
// 	if len(splitData) != 2 {
// 		return "", fmt.Errorf("nvalid base64 string")
// 	}

// 	mimeType := strings.Split(splitData[0], ";")[0]
// 	data := splitData[1]

// 	// Map the MIME type to a file extension
// 	mimeToExtension := map[string]string{
// 		"data:image/jpeg": ".jpg",
// 		"data:image/png":  ".png",
// 		"data:image/gif":  ".gif",
// 		// Add more mappings as needed
// 	}
// 	extension, ok := mimeToExtension[mimeType]
// 	if !ok {
// 		return "", fmt.Errorf("unsupported file type: %s", mimeType)
// 	}

// 	// Generate a UUID for the file name
// 	fileName := uuid.New().String() + extension

// 	// Decode the base64 string back to bytes
// 	decodedData, err := base64.StdEncoding.DecodeString(data)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Write the bytes to a new image file
// 	imagePath := filepath.Join(directoryPath, fileName)
// 	err = os.WriteFile(imagePath, decodedData, 0644)
// 	if err != nil {
// 		return "", err
// 	}

// 	return fileName, nil
// }
