package main

import (
	"TelegramFileUploader/types"
	"os"
)

// What is the file type to upload
var uploadType types.UploadFileType

// Where is the file
var filePath string

// Who should get the file
var receiverID = os.Getenv("RECEIVER_ID")

// The mime type of the file
var mimeType string
