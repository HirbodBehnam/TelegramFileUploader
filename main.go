package main

import (
	"TelegramFileUploader/metadata"
	"TelegramFileUploader/types"
	"context"
	"fmt"
	"github.com/gotd/td/telegram"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(1)
	// Parse the arguments
	if len(os.Args) == 1 {
		printUsage()
	}
	if len(os.Args) == 2 {
		filePath = os.Args[1]
		uploadType = types.UploadFileTypeDocument
	} else {
		filePath = os.Args[2]
		if err := uploadType.FromArgument(os.Args[1]); err != nil {
			fmt.Println("cannot parse the file type argument")
			printUsage()
		}
	}
	// Get the metadata
	var err error
	mimeType, err = metadata.GetMimeType(filePath)
	if err != nil {
		log.Fatalf("cannot get mime type of file: %s", err)
	}
	// Run the bot
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	err = telegram.BotFromEnvironment(ctx, telegram.Options{
		NoUpdates:      true, // don't subscribe to updates in one-shot mode
		SessionStorage: &telegram.FileSessionStorage{Path: "session"},
	}, nil, BeginBot)
	if err != nil {
		log.Fatalf("cannot upload: %s\n", err)
	}
}

func printUsage() {
	fmt.Println("Usage: uploader [-d|-p|-v|-a|-m] file_to_upload")
	fmt.Println("Flags:")
	fmt.Println("\t-d: Upload simple document")
	fmt.Println("\t-p: Upload as picture")
	fmt.Println("\t-v: Upload as video")
	fmt.Println("\t-a: Upload as voice")
	fmt.Println("\t-m: Upload as music")
	fmt.Println("Metadata of the file will be found out using ffmpeg and ffprobe")
	fmt.Println()
	fmt.Println("Environment variables which must be set: APP_ID APP_HASH BOT_TOKEN RECEIVER_ID")
	os.Exit(2)
}
