package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type YtdlMediaFormat struct {
	Asr               int         `json:"asr"`
	Filesize          int         `json:"filesize"`
	FormatID          string      `json:"format_id"`
	FormatNote        string      `json:"format_note"`
	Fps               interface{} `json:"fps"`
	Height            interface{} `json:"height"`
	Quality           int         `json:"quality"`
	Tbr               float64     `json:"tbr"`
	URL               string      `json:"url"`
	Width             interface{} `json:"width"`
	Ext               string      `json:"ext"`
	Vcodec            string      `json:"vcodec"`
	Acodec            string      `json:"acodec"`
	Abr               float64     `json:"abr,omitempty"`
	DownloaderOptions struct {
		HTTPChunkSize int `json:"http_chunk_size"`
	} `json:"downloader_options,omitempty"`
	Container   string `json:"container,omitempty"`
	Format      string `json:"format"`
	Protocol    string `json:"protocol"`
	HTTPHeaders struct {
		UserAgent      string `json:"User-Agent"`
		AcceptCharset  string `json:"Accept-Charset"`
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		AcceptLanguage string `json:"Accept-Language"`
	} `json:"http_headers"`
	Vbr float64 `json:"vbr,omitempty"`
}

type YtdlMetadata struct {
	ID         string            `json:"id"`
	Title      string            `json:"title"`
	Formats    []YtdlMediaFormat `json:"formats"`
	Thumbnails []struct {
		Height     int    `json:"height"`
		URL        string `json:"url"`
		Width      int    `json:"width"`
		Resolution string `json:"resolution"`
		ID         string `json:"id"`
	} `json:"thumbnails"`
	Description        string            `json:"description"`
	UploadDate         string            `json:"upload_date"`
	Uploader           string            `json:"uploader"`
	UploaderID         string            `json:"uploader_id"`
	UploaderURL        string            `json:"uploader_url"`
	ChannelID          string            `json:"channel_id"`
	ChannelURL         string            `json:"channel_url"`
	Duration           int               `json:"duration"`
	ViewCount          int               `json:"view_count"`
	AverageRating      interface{}       `json:"average_rating"`
	AgeLimit           int               `json:"age_limit"`
	WebpageURL         string            `json:"webpage_url"`
	Categories         []string          `json:"categories"`
	Tags               []string          `json:"tags"`
	IsLive             interface{}       `json:"is_live"`
	Channel            string            `json:"channel"`
	Extractor          string            `json:"extractor"`
	WebpageURLBasename string            `json:"webpage_url_basename"`
	ExtractorKey       string            `json:"extractor_key"`
	Playlist           interface{}       `json:"playlist"`
	PlaylistIndex      interface{}       `json:"playlist_index"`
	Thumbnail          string            `json:"thumbnail"`
	DisplayID          string            `json:"display_id"`
	RequestedSubtitles interface{}       `json:"requested_subtitles"`
	RequestedFormats   []YtdlMediaFormat `json:"requested_formats"`
	Format             string            `json:"format"`
	FormatID           string            `json:"format_id"`
	Width              int               `json:"width"`
	Height             int               `json:"height"`
	Resolution         interface{}       `json:"resolution"`
	Fps                int               `json:"fps"`
	Vcodec             string            `json:"vcodec"`
	Vbr                float64           `json:"vbr"`
	StretchedRatio     interface{}       `json:"stretched_ratio"`
	Acodec             string            `json:"acodec"`
	Abr                float64           `json:"abr"`
	Ext                string            `json:"ext"`
}

func IsYoutubeUrl(url string) bool {
	prefixes := []string{
		"https://www.youtu.be",
		"https://m.youtu.be",
		"https://youtu.be",
		"http://www.youtu.be",
		"http://m.youtu.be",
		"http://youtu.be",
		"https://www.youtube",
		"https://m.youtube",
		"https://youtube",
		"http://www.youtube",
		"http://m.youtube",
		"http://youtube",
	}

	for _, pre := range prefixes {
		if strings.HasPrefix(url, pre) {
			return true
		}
	}

	return false
}

func YoutubeMediaUrl(videoUrl string) (string, error) {
	cmd := exec.Command(
		"youtube-dl",
		"--dump-single-json",
		"--no-warnings",
		"--call-home",
		"--prefer-free-formats",
		"--youtube-skip-dash-manifest",
		"--rm-cache-dir",
		videoUrl,
	)

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var ytdlOutput YtdlMetadata

	if err := json.Unmarshal(stdout, &ytdlOutput); err != nil {
		log.Fatal(err)
		return "", err
	}

	// Just get the first available for now with opus
	for _, format := range ytdlOutput.Formats {
		if format.Vcodec == "none" && format.Acodec != "opus" {
			return format.URL, nil
		}
	}

	err = fmt.Errorf("No mediaurl found")
	log.Println(err)
	return "", err
}
