// Package main provides the entry point for the application.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aerosystems/slicer/internal/ffmpeg"
	"github.com/aerosystems/slicer/internal/stream"
	"golang.org/x/sync/errgroup"
)

func main() {
	ffmpegLoop, stdout, err := ffmpeg.NewLoopProcess("fixtures/bbb_sunflower_1080p_30fps_normal.mp4")
	if err != nil {
		log.Fatalf("error creating loop process: %v", err)
	}

	ffmpegHLS, stdin, err := ffmpeg.NewHLSProcess("hls")
	if err != nil {
		log.Fatalf("error creating hls process: %v", err)
	}

	streamService := stream.NewService(ffmpegLoop, stdout, ffmpegHLS, stdin)

	group := errgroup.Group{}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	group.Go(func() error {
		<-sig
		log.Println("received signal, stopping services")
		if err := ffmpegLoop.Process.Kill(); err != nil {
			log.Printf("error killing ffmpegLoop: %v", err)
		}
		if err := ffmpegHLS.Process.Kill(); err != nil {
			log.Printf("error killing ffmpegHLS: %v", err)
		}
		return nil
	})
	group.Go(func() error {
		return streamService.Run()
	})

	if err := group.Wait(); err != nil {
		log.Fatalf("error running stream service: %v", err)
	}
}
