package metadata

import (
	"context"
	"fmt"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
	"os"
	"path"
)

type FileHolder struct {
	File     tg.InputFileClass
	Mime     string
	Filepath string
}

// Document will make the file appear as simple document file
func (f FileHolder) Document() message.MediaOption {
	return message.UploadedDocument(f.File).
		MIME(f.Mime).
		ForceFile(true).
		Filename(path.Base(f.Filepath))
}

// Video will use ffprobe and ffmpeg to get length, resolution and thumbnail of video
// and then upload it
func (f FileHolder) Video(ctx context.Context, telegramUploader *uploader.Uploader) (message.MediaOption, error) {
	// At first get info from video
	width, height, duration, err := getVideoMetadata(ctx, f.Filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot get metadata of video: %w", err)
	}
	// Now generate a thumbnail
	thumbnailFile, err := f.thumbnailUploader(ctx, telegramUploader, generateVideoThumbnail)
	if err != nil {
		return nil, err
	}
	// Now create the file
	return message.UploadedDocument(f.File).
		MIME(f.Mime).
		Filename(path.Base(f.Filepath)).
		Thumb(thumbnailFile).
		Video().
		Resolution(width, height).
		DurationSeconds(duration).
		SupportsStreaming(), nil
}

// Photo will upload the file as photo
func (f FileHolder) Photo() message.MediaOption {
	return message.UploadedPhoto(f.File)
}

// Music will upload a music
func (f FileHolder) Music(ctx context.Context, telegramUploader *uploader.Uploader) (message.MediaOption, error) {
	// Create the document
	document := message.UploadedDocument(f.File).
		MIME(f.Mime).
		Filename(path.Base(f.Filepath))
	// Get metadata
	duration, title, performer, err := getAudioMetadata(ctx, f.Filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot get metadata: %w", err)
	}
	// Get thumbnail
	if thumbnail, err := f.thumbnailUploader(ctx, telegramUploader, generateAudioThumbnail); err == nil {
		document.Thumb(thumbnail)
	}
	// Create the music
	audio := document.Audio().DurationSeconds(duration)
	if title != "" {
		audio.Title(title)
	}
	if performer != "" {
		audio.Performer(performer)
	}
	return audio, nil
}

func (f FileHolder) Voice(ctx context.Context) (message.MediaOption, error) {
	// Get metadata again...
	duration, _, _, err := getAudioMetadata(ctx, f.Filepath)
	if err != nil {
		return nil, err
	}
	// Create the file
	return message.UploadedDocument(f.File).Filename(path.Base(f.Filepath)).
		Voice().DurationSeconds(duration), nil
}

func (f FileHolder) thumbnailUploader(ctx context.Context, telegramUploader *uploader.Uploader,
	generateFunction func(ctx context.Context, videoFile, thumbnail string) error) (tg.InputFileClass, error) {
	// At first create a temp file
	tempFile, err := os.CreateTemp("", "*.jpg")
	if err != nil {
		return nil, fmt.Errorf("cannot generate a temp file: %w", err)
	}
	_ = tempFile.Close()
	defer func(name string) {
		_ = os.Remove(name)
	}(tempFile.Name())
	// Now extract it
	err = generateFunction(ctx, f.Filepath, tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("cannot extract the thumbnail: %w", err)
	}
	// Upload it to telegram
	thumbnailFile, err := telegramUploader.FromPath(ctx, tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("cannot upload thumbnail: %w", err)
	}
	return thumbnailFile, nil
}
