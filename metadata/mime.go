package metadata

import "github.com/gabriel-vasile/mimetype"

// GetMimeType will get the mime type from a file
func GetMimeType(filename string) (string, error) {
	mime, err := mimetype.DetectFile(filename)
	if err != nil {
		return "", err
	}
	return mime.String(), nil
}
