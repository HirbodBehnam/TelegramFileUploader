package main

import (
	"TelegramFileUploader/metadata"
	"TelegramFileUploader/types"
	"context"
	"fmt"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

// BeginBot starts the bot and tries to upload the file to Telegram
func BeginBot(ctx context.Context, client *telegram.Client) error {
	api := tg.NewClient(client)
	u := uploader.NewUploader(api)
	sender := message.NewSender(api).WithUploader(u)
	// Upload the file
	upload, err := u.FromPath(ctx, filePath)
	if err != nil {
		return fmt.Errorf("upload %q: %w", filePath, err)
	}
	// Create the media
	err = nil
	file := metadata.FileHolder{
		File:     upload,
		Mime:     mimeType,
		Filepath: filePath,
	}
	var media message.MediaOption
	switch uploadType {
	case types.UploadFileTypeDocument:
		media = file.Document()
	case types.UploadFileTypeVideo:
		media, err = file.Video(ctx, u)
	case types.UploadFileTypePhoto:
		media = file.Photo()
	case types.UploadFileTypeMusic:
		media, err = file.Music(ctx, u)
	case types.UploadFileTypeVoice:
		if mimeType != message.DefaultVoiceMIME {
			err = fmt.Errorf("invalid mime type for ogg: expected %s got %s", message.DefaultVoiceMIME, mimeType)
			break
		}
		media, err = file.Voice(ctx)
	}
	if err != nil {
		return fmt.Errorf("cannot get metadata of files: %w", err)
	}
	// Send it
	_, err = sender.Resolve(receiverID).Media(ctx, media)
	// Done
	return err
}
