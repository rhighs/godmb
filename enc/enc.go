package enc

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"layeh.com/gopus"
)

type PlayerEvent int
type PlayerState int

const (
	PlayerEventTrackEnded PlayerEvent = iota
	PlayerEventPaused
	PlayerEventStopped
	PlayerEventResumed
)

const (
	PlayerStatePlaying PlayerState = iota
	PlayerStateStopped
	PlayerStatePaused
	PlayerStateIdle
)

type EncOptions struct {
	PcmOptions
	FrameSize     int
	Bitrate       int
	MaxCacheBytes int
}

type Enc struct {
	encoder   *gopus.Encoder
	listeners map[PlayerEvent][]func(PlayerEvent)
	State     PlayerState
}

type CacheOverflowError struct {
	MaxCacheBytes int
}

func (ps PlayerState) String() string {
	switch ps {
	case PlayerStatePlaying:
		return "PlayerStatePlaying"
	case PlayerStateStopped:
		return "PlayerStateStopped"
	case PlayerStatePaused:
		return "PlayerEventPaused"
	case PlayerStateIdle:
		return "PlayerStateIdle"
	}

	panic("unreachable!")
}

func (e *CacheOverflowError) Error() string {
	return "audio too large: the maximum cache limit of " + strconv.Itoa(e.MaxCacheBytes) + " bytes has been exceeded"
}

type Command interface{}
type CommandStop struct{}
type CommandPause struct{}
type CommandResume struct{}
type CommandSeek float32
type CommandGetPlaybackTime struct{}
type CommandGetDuration struct{} // Attempts to get the duration. Only succeeds if the encoder is already done.

type Response interface{}

type ResponsePlaybackTime float32     // Playback time in seconds.
type ResponseDuration float32         // Duration in seconds.
type ResponseDurationUnknown struct{} // Returned if the duration is unknown.

func DefaultOptions(ffmpegPath string) EncOptions {
	return EncOptions{
		PcmOptions:    getDefaultPcmOptions(ffmpegPath),
		FrameSize:     960,
		Bitrate:       64000,
		MaxCacheBytes: 4 * 1024 * 1024,
	}
}

func NewEnc(opts EncOptions) (out *Enc) {
	enc, err := gopus.NewEncoder(opts.SampleRate, opts.Channels, gopus.Audio)
	if err != nil {
		panic(err)
	}
	return &Enc{
		encoder:   enc,
		listeners: make(map[PlayerEvent][]func(PlayerEvent)),
		State:     PlayerStateIdle,
	}
}

func (e *Enc) Listen(event PlayerEvent, action func(PlayerEvent)) {
	if _, ok := e.listeners[event]; !ok {
		e.listeners[event] = make([]func(PlayerEvent), 0)
	}

	e.listeners[event] = append(e.listeners[event], action)
}

func (e *Enc) Notify(event PlayerEvent) {
	if actions, ok := e.listeners[event]; ok {
		for _, action := range actions {
			action(event)
		}
	}
}

func (e *Enc) GetOpusFrames(input string, opts EncOptions, ch chan<- []byte, errCh chan<- error, cmdCh <-chan Command, respCh chan<- Response) {
	pcm, cmd, err := getPcm(input, opts.PcmOptions)
	if err != nil {
		errCh <- err
		return
	}

	framesPerSecond := float32(opts.SampleRate) / float32(opts.FrameSize)
	maxSamples := opts.FrameSize * opts.Channels
	maxBytes := maxSamples * 2

	nof := 0

	sampleBytes := make([]byte, maxBytes)
	opusFrames := make([][]byte, 0, 512)
	samples := make([]int16, maxSamples)

	frameCh := make(chan []byte, 8)

	encoderDone := make(chan struct{})
	encoderStop := make(chan struct{})
	encoderResume := make(chan struct{})
	encoderPause := make(chan struct{})

	encoderRunning := true
	encoderPaused := false

	go func() {
		killedFfmpeg := false

	encoderLoop:
		for {
			select {
			case <-encoderStop:
				if err := cmd.Process.Signal(os.Interrupt); err != nil {
					fmt.Println("[ENCODER_ERR]: Failed sending interrupt for encoder to stop")
				}
				pcm.Close()
				killedFfmpeg = true
				break encoderLoop
			case <-encoderPause:
				encoderPaused = true
				<-encoderResume
				encoderPaused = false
			default:
			}

			_, err := io.ReadFull(pcm, sampleBytes)
			if err != nil {
				errCh <- err
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break
				}
			}

			// Fast way of reading bytes in couples, LE preserves order
			for i := range samples {
				samples[i] = int16(binary.LittleEndian.Uint16(sampleBytes[2*i:]))
			}

			// Encode to opus
			frame, err := e.encoder.Encode(samples, opts.FrameSize, maxBytes)
			if err != nil {
				errCh <- err
			}

			// Send frame to consumer
			frameCh <- frame
		}

		// Wait for ffmpeg to close.
		err = cmd.Wait()
		if err != nil {
			// Ffmpeg returns 255 if it was killed using SIGINT.
			switch e := err.(type) {
			case *exec.ExitError:
				if e.ExitCode() == 255 {
					if !killedFfmpeg {
						errCh <- err
					}
				}
			default:
				errCh <- err
			}
		}

		// Tell the main process that the encoder is done.
		encoderDone <- struct{}{}
	}()

	e.State = PlayerStatePlaying

	lastCacheSize := 0
	playerPaused := false

loop:
	for {
		select {
		case v := <-frameCh:
			opusFrames = append(opusFrames, v)

			cacheSize := 0
			for _, frames := range opusFrames[nof:] {
				cacheSize += len(frames)
			}

			if !encoderPaused && cacheSize >= opts.MaxCacheBytes {
				encoderPause <- struct{}{}
				opusFrames = opusFrames[nof:]
				nof = 0
			}

			lastCacheSize = cacheSize
		case <-encoderDone:
			encoderRunning = false
		case receivedCmd := <-cmdCh:
			switch v := receivedCmd.(type) {
			case CommandStop:
				e.State = PlayerStateStopped
				e.Notify(PlayerEventStopped)
				if encoderRunning {
					encoderStop <- struct{}{}
				}
				break loop
			case CommandPause:
				e.State = PlayerStatePaused
				playerPaused = true
				e.Notify(PlayerEventPaused)
			case CommandResume:
				e.State = PlayerStatePlaying
				playerPaused = false
				e.Notify(PlayerEventResumed)
			case CommandSeek:
				nof = int(float32(v) * framesPerSecond)
			case CommandGetPlaybackTime:
				respCh <- ResponsePlaybackTime(float32(nof / int(framesPerSecond)))
			case CommandGetDuration:
				respCh <- ResponseDuration(float32(len(opusFrames)) / framesPerSecond)
			}
		default:
			time.Sleep(2 * time.Millisecond)
		}

		if encoderPaused && !playerPaused {
			cacheSizeLeft := 0
			for _, frames := range opusFrames[nof:] {
				cacheSizeLeft += len(frames)
			}

			if cacheSizeLeft < lastCacheSize/3 {
				encoderResume <- struct{}{}
			}
		}

		if !playerPaused && nof < len(opusFrames) {
			select {
			case ch <- opusFrames[nof]:
				nof++
			default:
			}
		}

		if !encoderRunning && nof >= len(opusFrames) {
			break
		}
	}

	close(respCh)
	e.State = PlayerStateIdle
	e.Notify(PlayerEventTrackEnded)
}
