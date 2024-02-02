package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pandodao/tokenizer-go"
)

var TokenLimit = 4000

type PostType int

const (
	PostTypeNews PostType = iota
	PostTypeColumn
	PostTypeInterview
	PostTypeReport
	PostTypeEvent
	PostTypeReview
)

func (p PostType) String() string {
	switch p {
	case PostTypeNews:
		return "뉴스"
	case PostTypeColumn:
		return "칼럼"
	case PostTypeInterview:
		return "인터뷰"
	case PostTypeReport:
		return "리포트"
	case PostTypeEvent:
		return "이벤트"
	case PostTypeReview:
		return "리뷰"
	default:
		return "알 수 없음"
	}
}

type Post struct {
	Type       PostType `json:"type"`
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	Categories []string `json:"categories"`
	Image      string   `json:"image"`
	Contents   string   `json:"contents"`
	Summary    string   `json:"summary"`
	Summarized bool     `json:"summarized,omitempty"`
}

func (p Post) String() string {
	return fmt.Sprintf("제목: %s\n카테고리: %s\nURL: %s\n이미지: %s\n내용:\n%s\n요약:\n%s\n", p.Title, strings.Join(p.Categories, ", "), p.URL, p.Image, p.Contents, p.Summary)
}

func (p Post) FormatSummarizable() (string, bool) {
	s := fmt.Sprintf("제목: %s\n카테고리: %s\n내용:\n%s\n", p.Title, strings.Join(p.Categories, ", "), p.Contents)

	tokens, err := tokenizer.CalToken(s)
	if err != nil {
		return "", false
	}

	if tokens > TokenLimit {
		return "", false
	}

	return s, true
}

func (p Post) ToMessage() *Message {
	var embed MessageEmbed

	if p.Summarized {
		embed = MessageEmbed{
			Title:       p.Title,
			Description: p.Summary,
			URL:         p.URL,
			Image: &MessageEmbedImage{
				URL: p.Image,
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
	} else {
		embed = MessageEmbed{
			Title:       p.Title,
			Description: p.Contents, // TODO: it might be too long
			URL:         p.URL,
			Image: &MessageEmbedImage{
				URL: p.Image,
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	return &Message{
		Embeds: []*MessageEmbed{&embed},
	}
}
