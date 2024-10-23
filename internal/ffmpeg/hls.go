// Package ffmpeg provides a simple interface to the ffmpeg command line tool.
package ffmpeg

import (
	"io"
	"os/exec"
)

// NewHLSProcess creates a new ffmpeg process that generates an HLS stream.
func NewHLSProcess(dstFolder string) (*exec.Cmd, io.WriteCloser, error) {
	cmd := exec.Command("ffmpeg", //nolint: gosec
		"-i", "pipe:0",
		"-codec:v", "copy",
		"-codec:a", "aac",
		"-f", "hls",
		"-hls_time", "6",
		"-hls_list_size", "8",
		`-strftime`, "1",
		"-hls_segment_filename", dstFolder+"/%s.ts",
		dstFolder+"/index.m3u8")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	return cmd, stdin, nil
}
