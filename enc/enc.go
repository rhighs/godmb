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
	encoderResume := make(chan struct{})
	encoderPause := make(chan struct{})

	encoderRunning := true
	paused := false
	loop := false

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
				<-encoderResume // Wait for a resume signal
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

			if cacheSize > opts.MaxCacheBytes {
				errCh <- &CacheOverflowError{
					MaxCacheBytes: opts.MaxCacheBytes,
				}
				break
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
				encoderPause <- struct{}{}
				e.Notify(PlayerEventPaused)
			case CommandResume:
				paused = false
				e.State = PlayerStatePlaying
				encoderResume <- struct{}{}
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
				respCh <- ResponseDuration(float32(len(opusFrames)) / framesPerSecond)
			}
		default:
			time.Sleep(2 * time.Millisecond)
		}

        // Audio is playing, send it to caller channel
		if !paused && nof < len(opusFrames) {
			select {
			case ch <- opusFrames[nof]:
				cacheSize -= len(opusFrames[nof])
				nof++
			default:
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

		// Cut the cache if we're going too far with the memory
		if len(opusFrames) >= opts.MaxCacheBytes {
			opusFrames = opusFrames[nof:]
			nof = 0
			cacheSize = len(opusFrames)
		}
	}

	// Wait for the encoder to finish if it's still running.
	if encoderRunning {
		<-encoderDone
		close(respCh)
	}

	e.State = PlayerStateIdle
	e.Notify(PlayerEventTrackEnded)
}
