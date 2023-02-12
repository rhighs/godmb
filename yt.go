package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

const (
	REGEXP_PREFIX string = `<meta name="title" content="`
	LOG_YTCMD     bool   = true
)

var YtdlpPath string = ""

func SetYtdlpPath(path string) {
	YtdlpPath = path
	c := exec.Command(YtdlpPath)
	if err := c.Start(); err != nil {
		log.Fatalf("[COMMAND_ERR] yt-dlp start failed with reason: %v", err)
	}

	if err := c.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 2 {
			// Exit code 2 means an argument-less command invocation
			log.Fatalf("[COMMAND_ERR] yt-dlp exited with code %d Error: %s", exitErr.ExitCode(), exitErr.Error())
		}
	}
}

func GetYtdlpPath() string {
	return YtdlpPath
}

var ytPrefixes = []string{
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

var videoTitleRegexp *regexp.Regexp

func init() {
	sr := `<meta name="title" content=".*">`
	r, err := regexp.Compile(sr)
	if err != nil {
		panic("Couldn't compile regexp: " + sr)
	}
	videoTitleRegexp = r
}

func IsYoutubeUrl(url string) bool {
	for _, pre := range ytPrefixes {
		if strings.HasPrefix(url, pre) {
			return true
		}
	}
	return false
}

func GetYTVideoTitle(URL string) (string, error) {
	resp, _ := http.Get(URL)
	bodyBytes, _ := io.ReadAll(resp.Body)
	matched := videoTitleRegexp.FindAllStringSubmatch(string(bodyBytes), 1)
	if len(matched[0]) == 0 {
		return "", errors.New("Could find any title at: " + URL)
	}
	str := matched[0][0]
	idx := strings.Index(str, ">")
	str = str[len(REGEXP_PREFIX) : idx-1]
	return str, nil
}

func YoutubeMediaUrl(videoUrl string) (string, error) {
	args := []string{
		"--dump-single-json",
		"--no-warnings",
		"--call-home",
		"--prefer-free-formats",
		"--youtube-skip-dash-manifest",
		"--rm-cache-dir",
		videoUrl,
	}

	cmd := exec.Command(
		YtdlpPath,
		args...,
	)

	if LOG_YTCMD {
		log.Println("[YTDL_CMD_USED]:", YtdlpPath, strings.Join(args, " "))
	}

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var ytdlOutput YTDLPOut

	if err := json.Unmarshal(stdout, &ytdlOutput); err != nil {
		log.Fatal(err)
		return "", err
	}

	// Just get the first one including opus
	for _, format := range ytdlOutput.Formats {
		if format.Vcodec == "none" && format.Acodec == "opus" {
			return format.URL, nil
		}
	}

	err = fmt.Errorf("no media url found")
	log.Println(err)
	return "", err
}
