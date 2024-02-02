package main

import (
	"os"
	"syscall"
	"time"
)

func main() {
	// init logger
	l, _ := newLogger(false, "service", "gorypto news")
	defer l.Sync()

	// init scraper
	ts := NewTokenPostScraper(true)

	// TODO: init summarizer

	// init scheduler
	s, err := NewScheduler(nil)
	if err != nil {
		l.Fatal("Failed to create scheduler", "error", err)
	}

	// add scraper to scheduler
	summarizedPosts := make(chan *Post)

	err = s.AddScraper(ts, summarizedPosts, 1*time.Hour, 1, true)
	if err != nil {
		l.Fatal("Failed to add scraper to scheduler", "error", err)
	}

	// init webhook
	w := NewDiscordWebhook(nil, os.Getenv("WEBHOOK_URL"), true)

	// run webhook and scheduler
	w.Run()
	s.Start()

	// process summarized posts to webhook message
	// and send it via webhook
	go func() {
		for p := range summarizedPosts {
			msg := &Message{
				Content: p.Contents,
			}

			w.Send(msg)
		}
	}()

	// graceful shutdown on SIGINT and SIGTERM
	<-GracefulShutdown(func() {
		close(summarizedPosts)
		w.Stop()
		s.StopJobs()
	}, syscall.SIGINT, syscall.SIGTERM)
}
