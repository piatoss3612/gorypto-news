package main

import (
	"os"
	"syscall"
)

func main() {
	l, _ := newLogger(false, "service", "gorypto news")
	defer l.Sync()

	// s := NewTokenPostScraper(true)

	// posts, done, errs := s.Scrape(1)

	// for {
	// 	select {
	// 	case post := <-posts:
	// 		if post == nil {
	// 			continue
	// 		}
	// 		fmt.Println(post.String())
	// 	case <-done:
	// 		fmt.Println("Done")
	// 		return
	// 	case err := <-errs:
	// 		fmt.Println(err)
	// 	}
	// }

	w := NewDiscordWebhook(nil, os.Getenv("WEBHOOK_URL"), true)

	w.Run()

	msg := &Message{
		Content: "Hello, world!",
	}

	w.Send(msg)

	<-GracefulShutdown(func() {
		w.Stop()
	}, syscall.SIGINT, syscall.SIGTERM)
}
