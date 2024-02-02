package main

import (
	"flag"
	"os"
	"syscall"
	"time"

	"github.com/tmc/langchaingo/llms/openai"
	"github.com/valyala/fasthttp"
)

func main() {
	limit := flag.Int("limit", 5, "limit of posts to scrape")
	flag.Parse()

	if *limit < 1 {
		*limit = 1
	}

	// init logger
	l, _ := newLogger(false, "service", "gorypto news")
	defer l.Sync()

	// init scraper
	ts := NewTokenPostScraper(true)

	// init summarizer
	llm, err := openai.New()
	if err != nil {
		l.Fatal("Failed to create OpenAI LLM", "error", err)
	}

	cache := NewInMemoryCache()

	sum := NewSummarizer(llm, cache)

	// init scheduler
	s, err := NewScheduler(sum)
	if err != nil {
		l.Fatal("Failed to create scheduler", "error", err)
	}

	// add scraper to scheduler
	summarizedPosts := make(chan *Post)

	err = s.AddScraper(ts, summarizedPosts, 1*time.Hour, uint(*limit), true)
	if err != nil {
		l.Fatal("Failed to add scraper to scheduler", "error", err)
	}

	// init webhook
	client := fasthttp.Client{
		Name: "gorypto-news",
	}

	w := NewDiscordWebhook(&client, os.Getenv("WEBHOOK_URL"), true)

	// run webhook and scheduler
	l.Info("Running scheduler")

	w.Run()
	s.Start()

	// process summarized posts to webhook message
	// and send it via webhook
	go func() {
		for p := range summarizedPosts {
			if p == nil {
				continue
			}

			l.Info("Sending post via webhook", "id", p.ID)

			w.Send(p.ToMessage())
		}
	}()

	// graceful shutdown on SIGINT and SIGTERM
	<-GracefulShutdown(func() {
		l.Info("Shutting down...")
		close(summarizedPosts)
		w.Stop()
		s.StopJobs()
	}, syscall.SIGINT, syscall.SIGTERM)

	l.Info("Shutdown complete")
}
