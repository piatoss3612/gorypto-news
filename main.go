package main

import "fmt"

func main() {
	l, _ := newLogger(false, "service", "gorypto news")
	defer l.Sync()

	s := NewTokenPostScraper(true)

	posts, done, errs := s.Scrape(5)

	for {
		select {
		case post := <-posts:
			if post == nil {
				continue
			}
			fmt.Println(post.String())
		case <-done:
			fmt.Println("Done")
			return
		case err := <-errs:
			fmt.Println(err)
		}
	}
}
