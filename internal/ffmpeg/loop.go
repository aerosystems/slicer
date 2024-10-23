// Package ffmpeg provides a simple interface to the ffmpeg command line tool.
package ffmpeg

import (
	"io"
	"os/exec"
)

// NewLoopProcess creates a new ffmpeg process that loops the input file.
func NewLoopProcess(inputFile string) (*exec.Cmd, io.ReadCloser, error) {
	cmd := exec.Command("ffmpeg",
		"-re",
		"-stream_loop", "-1",
		"-i", inputFile,
		"-f", "mpegts",
		"-codec:v", "mpeg1video",
		"-")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	return cmd, stdout, nil
}
