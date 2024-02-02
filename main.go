package main

import "fmt"

func main() {
	_, _ = newLogger(true, "service", "gorypto news")

	s := NewTokenPostScraper(true)

	posts, done, errs := s.Scrape()

	for {
		select {
		case post := <-posts:
			fmt.Println(post.String())
		case <-done:
			fmt.Println("Done")
			return
		case err := <-errs:
			fmt.Println(err)
		}
	}
}
