package metadata

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

// getVideoMetadata will get simple metadata from video
func getVideoMetadata(ctx context.Context, filename string) (width, height, duration int, err error) {
	// Run ffprobe
	cmd := exec.CommandContext(ctx, "ffprobe", "-i", filename, "-show_format",
		"-show_entries", "stream=width,height")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, err
	}
	// Regex out the metadata
	width, err = matchAndParseRegex(out, "width=([0-9]+)")
	if err != nil {
		return 0, 0, 0, err
	}
	height, err = matchAndParseRegex(out, "height=([0-9]+)")
	if err != nil {
		return 0, 0, 0, err
	}
	duration, err = matchAndParseRegex(out, "duration=([0-9]+)")
	return
}

// generateVideoThumbnail will generate a thumbnail from a video file
func generateVideoThumbnail(ctx context.Context, videoFile, thumbnail string) error {
	ffmpegCmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", videoFile, "-vf", "scale=-2:360",
		"-vframes", "1", thumbnail)
	stderr := new(bytes.Buffer)
	ffmpegCmd.Stderr = stderr
	err := ffmpegCmd.Run()
	if err != nil {
		log.Println("trace:", stderr.String())
		return fmt.Errorf("cannot generate thumbnail: %w", err)
	}
	return nil
}

// generateAudioThumbnail will generate a thumbnail from an audio file
// Will return an error if there is no thumbnail
func generateAudioThumbnail(ctx context.Context, audioFile, thumbnail string) error {
	ffmpegCmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", audioFile, "-vf", "scale=-2:360", thumbnail)
	err := ffmpegCmd.Run()
	if err != nil {
		return fmt.Errorf("cannot generate thumbnail: %w", err)
	}
	return nil
}

func getAudioMetadata(ctx context.Context, filename string) (duration int, title, performer string, err error) {
	// Run ffprobe
	cmd := exec.CommandContext(ctx, "ffprobe", "-i", filename, "-show_format")
	out, err := cmd.Output()
	if err != nil {
		return 0, "", "", err
	}
	// Get the duration
	duration, err = matchAndParseRegex(out, "duration=([0-9]+)")
	if err != nil {
		return 0, "", "", err
	}
	// Get the tags
	title = matchRegex(out, "TAG:title=(.+)")
	performer = matchRegex(out, "TAG:artist=(.+)")
	return
}

// matchAndParseRegex will match a byte array against a pattern and extract an int from it
func matchAndParseRegex(data []byte, pattern string) (int, error) {
	regex := regexp.MustCompile(pattern)
	regexMatch := regex.FindAllSubmatch(data, 1)
	if len(regexMatch) == 0 {
		return 0, errors.New("no match found")
	}
	return strconv.Atoi(string(regexMatch[0][1]))
}

func matchRegex(data []byte, pattern string) string {
	regex := regexp.MustCompile(pattern)
	regexMatch := regex.FindAllSubmatch(data, 1)
	if len(regexMatch) == 0 {
		return ""
	}
	return string(regexMatch[0][1])
}
