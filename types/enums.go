package types

import "errors"

// UploadFileType determines the type of file to upload
type UploadFileType uint8

const (
	UploadFileTypeDocument UploadFileType = iota
	UploadFileTypePhoto
	UploadFileTypeVideo
	UploadFileTypeVoice
	UploadFileTypeMusic
)

// FromArgument converts a command line argument to UploadFileType and returns an error
// if the argument is invalid
func (u *UploadFileType) FromArgument(argument string) error {
	switch argument {
	case "-d":
		*u = UploadFileTypeDocument
	case "-p":
		*u = UploadFileTypePhoto
	case "-v":
		*u = UploadFileTypeVideo
	case "-a":
		*u = UploadFileTypeVoice
	case "-m":
		*u = UploadFileTypeMusic
	default:
		return errors.New("invalid argument")
	}
	return nil
}
