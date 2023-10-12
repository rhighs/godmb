package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"ndmb/enc"

	"github.com/Pauloo27/searchtube"
	dgo "github.com/bwmarrin/discordgo"
)

var FfmpegPath string = ""

func SetFfmpegPath(path string) {
	FfmpegPath = path
}

func GetFfmpegPath() string {
	return FfmpegPath
}

func ResolveAudioSource(input string) (Track, error) {
	webUrl := ""
	track := Track{}

	if IsYoutubeUrl(input) {
		mediaUrl, err := YoutubeMediaUrl(input)
		if err != nil {
			return track, err
		}

		track.WebURL = input
		track.MediaURL = mediaUrl
		track.Title, err = ResolveVideoTitle(input)

		if err != nil {
			return track, err
		}

		return track, nil
	}

	// Might be a direct http stream
	if strings.HasPrefix(input, "http") {
		track.Title = "Unknown"
		track.MediaURL = input
		track.WebURL = "Unknown source"
		return track, nil
	}

	results := []*searchtube.SearchResult{}
	searchResults, err := searchtube.Search(input, 1)
	if err != nil {
		return track, err
	}

	for _, r := range searchResults {
		if !r.Live {
			results = append(results, r)
		}
	}

	webUrl = results[0].URL

	mediaUrl, err := YoutubeMediaUrl(webUrl)
	if err != nil {
		return track, err
	}

	track.MediaURL = mediaUrl
	track.WebURL = webUrl
	track.Title, err = ResolveVideoTitle(webUrl)
	if err != nil {
		return track, err
	}

	return track, nil
}

func FetchHttpMediaStream(mediaUrl string) (io.ReadCloser, int64, error) {
	req, err := http.NewRequest("GET", mediaUrl, nil)
	if err != nil {
		log.Println("Error preparing media stream request:", err)
		return nil, 0, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending media stream request:", err)
		return nil, 0, err
	}

	return resp.Body, resp.ContentLength, nil
}

func PlayMediaInVoiceChannel(mediaUrl string, player *enc.Enc, voiceConnection *dgo.VoiceConnection, errCh chan error, cmdCh chan enc.Command, respCh chan enc.Response) {
	go player.GetOpusFrames(mediaUrl, enc.DefaultOptions(FfmpegPath), voiceConnection.OpusSend, errCh, cmdCh, respCh)
}
