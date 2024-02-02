package main

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	"github.com/gocolly/colly/v2"
)

const (
	TokenPostBaseURL  = "https://www.tokenpost.kr"
	TokenPostCacheDir = "./tokenpost_cache"
)

type Scraper interface {
	Scrape(limit uint) (<-chan *Post, <-chan struct{}, <-chan error)
	Close() error
}

type TokenPostScraper struct {
	*colly.Collector

	l       *Logger
	logging bool
}

func NewTokenPostScraper(logging bool) *TokenPostScraper {
	s := &TokenPostScraper{
		Collector: colly.NewCollector(
			colly.AllowedDomains("www.tokenpost.kr", "tokenpost.kr"),
			// colly.CacheDir(TokenPostCacheDir), // arising error in docker container (permission denied)
			colly.Async(),
		),
		logging: logging,
	}

	if logging {
		s.l = GetLogger()
	}

	return s
}

func (s *TokenPostScraper) Scrape(limit uint) (<-chan *Post, <-chan struct{}, <-chan error) {
	if limit == 0 {
		return nil, nil, nil
	}

	c := s.Collector

	count := atomic.Value{}
	count.Store(uint(0))

	errs := make(chan error)
	posts := make(chan *Post)

	detailCollector := c.Clone()

	c.OnRequest(func(r *colly.Request) {
		if s.logging {
			s.l.Debug("Visiting", "URL", r.URL.String())
		}
	})

	c.OnHTML(`div[id=content] div.list_item_title`, func(e *colly.HTMLElement) {
		cur := count.Load().(uint)

		if cur >= limit {
			return
		}

		postURL := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		if postURL == "" {
			return
		}

		count.Store(cur + 1)

		detailCollector.Visit(postURL)
	})

	c.OnScraped(func(r *colly.Response) {
		if s.logging {
			s.l.Debug("Finished", "URL", r.Request.URL.String())
		}
	})

	detailCollector.OnRequest(func(r *colly.Request) {
		if s.logging {
			s.l.Debug("Visiting", "URL", r.URL.String())
		}
	})

	detailCollector.OnHTML(`div[id=content] div[id=articleContentArea]`, func(e *colly.HTMLElement) {
		categories := e.ChildTexts("div.view_blockchain_item > span")
		title := strings.TrimSpace(e.ChildText("span.view_top_title"))
		img := strings.TrimSpace(e.ChildAttr("div.imgBox > img", "src"))

		builder := strings.Builder{}

		e.ForEach("div.article_content > p", func(_ int, h *colly.HTMLElement) {
			content := strings.TrimSpace(h.Text)
			if content == "" {
				return
			}

			if strings.Contains(content, "[email") {
				return
			}

			strongs := h.ChildTexts("strong")

			if len(strongs) > 0 {
				for _, strong := range strongs {
					content = strings.ReplaceAll(content, strong, fmt.Sprintf("**%s**", strong))
				}
			}

			builder.WriteString(fmt.Sprintf("%s\n\n", content))
		})

		posts <- &Post{
			Type:       PostTypeNews,
			ID:         fmt.Sprintf("TokenPost%s", e.Request.URL.Path),
			Title:      title,
			Categories: categories,
			URL:        e.Request.URL.String(),
			Image:      img,
			Contents:   builder.String(),
		}
	})

	detailCollector.OnScraped(func(r *colly.Response) {
		if s.logging {
			s.l.Debug("Finished", "URL", r.Request.URL.String())
		}
	})

	done := make(chan struct{})

	go func() {
		c.Visit("https://www.tokenpost.kr/blockchain")
		c.Wait()
		detailCollector.Wait()

		close(done)
		close(posts)
		close(errs)
	}()

	return posts, done, errs
}

func (s *TokenPostScraper) ClearCache() error {
	return os.RemoveAll(TokenPostCacheDir)
}

func (s *TokenPostScraper) Close() error {
	return s.ClearCache()
}
