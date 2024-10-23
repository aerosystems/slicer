// Package stream service for processing streams.
package stream

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

// Service provides a stream service.
type Service struct {
	ffmpegLoop       *exec.Cmd
	ffmpegLoopOutput io.ReadCloser
	ffmpegHLS        *exec.Cmd
	ffmpegHLSInput   io.WriteCloser
}

// NewService creates a new stream service.
func NewService(ffmpegLoop *exec.Cmd, ffmpegLoopOutput io.ReadCloser, ffmpegHLS *exec.Cmd, ffmpegHLSInput io.WriteCloser) *Service {
	return &Service{
		ffmpegLoop:       ffmpegLoop,
		ffmpegLoopOutput: ffmpegLoopOutput,
		ffmpegHLS:        ffmpegHLS,
		ffmpegHLSInput:   ffmpegHLSInput,
	}
}

// Run starts the stream service.
func (s *Service) Run() error {
	go func() {
		defer s.ffmpegHLSInput.Close() //nolint: errcheck
		tee := io.TeeReader(s.ffmpegLoopOutput, s.ffmpegHLSInput)
		buf := make([]byte, 1024) //nolint:mnd
		for {
			n, err := tee.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatalf("error reading from tee: %v", err)
			}
			if n == 0 {
				break
			}

			// log.Printf("err: %v, n: %d, buf: %s", err, n, buf[:n])
			log.Printf("read %d bytes", n)
		}
	}()

	if err := s.ffmpegLoop.Wait(); err != nil {
		return fmt.Errorf("error waiting for ffmpegLoop: %w", err)
	}

	return s.ffmpegHLS.Wait()
}
