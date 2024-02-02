package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pandodao/tokenizer-go"
	"github.com/tmc/langchaingo/llms/openai"
)

var (
	ErrInvalidPost   = fmt.Errorf("invalid post")
	ErrTooManyTokens = fmt.Errorf("too many tokens")
)

var prompt = "당신은 블록체인과 관련된 전문적인 지식을 갖추고 있습니다. 아래의 블록체인 관련 게시글의 본문에는 마크다운 형식(**제목**)으로 소제목이 포함되어 있을 수 있습니다. 이를 요약해주세요."

type Summarizer struct {
	llm   *openai.LLM
	cache Cache
}

func NewSummarizer(llm *openai.LLM, cache Cache) *Summarizer {
	return &Summarizer{
		llm:   llm,
		cache: cache,
	}
}

func (s *Summarizer) Summarize(ctx context.Context, post *Post) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if post == nil {
		return ErrInvalidPost
	}

	if s.cache.Exists(ctx, post.ID) {
		var summary string
		err := s.cache.Get(ctx, post.ID, &summary)
		if err != nil {
			return err
		}

		post.Summary = summary
		post.Summarized = true

		return nil
	}

	content := post.FormatSummarizable()

	token, err := tokenizer.CalToken(content)
	if err != nil {
		return err
	}

	if token > 4500 {
		return ErrTooManyTokens
	}

	resp, err := s.llm.Call(ctx, prompt)
	if err != nil {
		return err
	}

	summary := resp

	err = s.cache.Set(ctx, post.ID, summary, time.Hour*24*3)
	if err != nil {
		return err
	}

	return nil
}
