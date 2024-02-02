package main

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	gocron.Scheduler
	sum *Summarizer

	l *Logger
}

func NewScheduler(sum *Summarizer) (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	scheduler := &Scheduler{
		Scheduler: s,
		sum:       sum,
		l:         GetLogger(),
	}

	return scheduler, nil
}

func (s *Scheduler) AddScraper(scraper Scraper, res chan<- *Post, duration time.Duration, limit uint, logging bool) error {
	_, err := s.NewJob(gocron.DurationJob(duration), gocron.NewTask(func() {
		post, done, errs := scraper.Scrape(limit)

		for {
			select {
			case p := <-post:
				if p == nil {
					continue
				}

				// err := s.summarizer.Summarize(context.Background(), p)
				// if err != nil {
				// 	if logging {
				// 		Error("Failed to summarize post", err)
				// 	}

				// 	continue
				// }

				res <- p
			case <-done:
				if logging {
					s.l.Debug("Done scraping posts")
				}

				return
			case err := <-errs:
				if err == nil {
					continue
				}

				if logging {
					s.l.Error("Failed to scrape post", "error", err)
				}
			}
		}
	}))
	if err != nil {
		return err
	}

	return nil
}
