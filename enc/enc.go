package enc

/*
Most of the credit for this opus buffering stuff here goes to the guy who made this publicly available:
https://github.com/xypwn/go-musicbot/blob/master/dca0/dca0.go

The encoder is fairly simple and supports a variety of commands. Everything involving PCM is intended as processes for receiving data,
this data comes from ffpmeg which is then fed to `gopus` for encoding. In between, there is a buffering mechanism
made possible by the use of goroutines, therefore you'll see the encoder buffering onto a `frameCh` channel.
*/

import (
	"encoding/binary"
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

func (e *CacheOverflowError) Error() string {
	return "audio too large: the maximum cache limit of " + strconv.Itoa(e.MaxCacheBytes) + " bytes has been exceeded"
}

type Command interface{}
type CommandStop struct{}
type CommandPause struct{}
type CommandResume struct{}
type CommandStartLooping struct{}
type CommandStopLooping struct{}
type CommandSeek float32             // In seconds.
type CommandGetPlaybackTime struct{} // Gets the playback time.
type CommandGetDuration struct{}     // Attempts to get the duration. Only succeeds if the encoder is already done.

type Response interface{}

type ResponsePlaybackTime float32     // Playback time in seconds.
type ResponseDuration float32         // Duration in seconds.
type ResponseDurationUnknown struct{} // Returned if the duration is unknown.

func DefaultOptions(ffmpegPath string) EncOptions {
	return EncOptions{
		PcmOptions:    getDefaultPcmOptions(ffmpegPath),
		FrameSize:     960,
		Bitrate:       64000,
		MaxCacheBytes: 20000000,
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

	cacheSize := 0
	opusFrames := make([][]byte, 0, 512)
	nof := 0
	frameCh := make(chan []byte, 8)
	sampleBytes := make([]byte, maxBytes)
	samples := make([]int16, maxSamples)

	encoderDone := make(chan struct{})
	encoderStop := make(chan struct{})

	go func() {
		killedFfmpeg := false
	encoderLoop:
		for {
			select {
			case <-encoderStop:
				cmd.Process.Signal(os.Interrupt)
				pcm.Close()
				killedFfmpeg = true
				break encoderLoop
			default:
			}

			_, err := io.ReadFull(pcm, sampleBytes)
			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break
				} else {
					errCh <- err
				}
			}

			// Fast way of reading bytes in couples, LE preserves order
			for i := range samples {
				samples[i] = int16(binary.LittleEndian.Uint16(sampleBytes[2*i:]))
			}

			// Encode to opus, prepare sending (or other stuff)
			frame, err := e.encoder.Encode(samples, opts.FrameSize, maxBytes)
			if err != nil {
				errCh <- err
			}

			if cacheSize > opts.MaxCacheBytes {
				errCh <- &CacheOverflowError{
					MaxCacheBytes: opts.MaxCacheBytes,
				}
				break
			}

			cacheSize += len(frame)
			frameCh <- frame
		}

		// Wait for ffmpeg to close.
		err = cmd.Wait()
		if err != nil {
			// Ffmpeg returns 255 if it was killed using SIGINT. Therefore that
			// wouldn't be an error.
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

	encoderRunning := true
	paused := false
	loop := false

	e.State = PlayerStatePlaying

loop:
	for {
		select {
		case v := <-frameCh:
			opusFrames = append(opusFrames, v)
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
				paused = true
				e.State = PlayerStatePaused
				e.Notify(PlayerEventPaused)
			case CommandResume:
				paused = false
				e.State = PlayerStatePlaying
				e.Notify(PlayerEventResumed)
			case CommandStartLooping:
				e.State = PlayerStatePlaying
				loop = true
			case CommandSeek:
				nof = int(float32(v) * framesPerSecond)
			case CommandStopLooping:
				loop = false
			case CommandGetPlaybackTime:
				respCh <- ResponsePlaybackTime(float32(nof / int(framesPerSecond)))
			case CommandGetDuration:
				//if encoderRunning {
				//	respCh <- ResponseDurationUnknown{}
				//} else {
				//	respCh <- ResponseDuration(float32(len(opusFrames)) / framesPerSecond)
				//}
				respCh <- ResponseDuration(float32(len(opusFrames)) / framesPerSecond)
			}
		default:
			time.Sleep(2 * time.Millisecond)
		}

		if !paused && nof < len(opusFrames) {
			if encoderRunning {
				select {
				case ch <- opusFrames[nof]:
					nof++
				default:
				}
			} else {
				ch <- opusFrames[nof]
				nof++
			}
		}

		if !encoderRunning && nof >= len(opusFrames) {
			if loop {
				nof = 0
			} else {
				// We're done sending opus data.
				break
			}
		}
	}

	// Wait for the encoder to finish if it's still running.
	if encoderRunning {
		<-encoderDone
		encoderRunning = false
		// TODO: I want to make this unnecessary. I have just noticed that
		// this channel often get stuck so a panic is more helpful than that.
		close(respCh)
	}

	e.State = PlayerStateIdle
	e.Notify(PlayerEventTrackEnded)
}
