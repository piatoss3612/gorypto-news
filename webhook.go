package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// types from https://github.com/bwmarrin/discordgo
type Message struct {
	// The content of the message.
	Content string `json:"content"`

	// A list of embeds present in the message.
	Embeds []*MessageEmbed `json:"embeds"`
}

type MessageEmbed struct {
	URL         string                 `json:"url,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []*MessageEmbedField   `json:"fields,omitempty"`
}

type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type MessageEmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type MessageEmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type MessageEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type MessageEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type DiscordWebhook struct {
	client *fasthttp.Client
	URL    string

	l       *Logger
	logging bool

	msgChan  chan *Message
	doneChan chan struct{}
}

func NewDiscordWebhook(client *fasthttp.Client, URL string, logging bool) *DiscordWebhook {
	w := &DiscordWebhook{
		URL:      URL,
		client:   client,
		l:        GetLogger(),
		logging:  logging,
		msgChan:  make(chan *Message, 10),
		doneChan: make(chan struct{}),
	}

	if client == nil {
		w.client = &fasthttp.Client{}
	}

	return w
}

func (w *DiscordWebhook) Run() {
	go w.handle()
}

func (w *DiscordWebhook) Stop() {
	close(w.doneChan)
}

func (w *DiscordWebhook) Send(msg *Message) {
	w.msgChan <- msg
}

func (w *DiscordWebhook) handle() {
	defer close(w.msgChan)

	for {
		select {
		case msg := <-w.msgChan:
			if msg == nil {
				continue
			}

			time.Sleep(1 * time.Second) //

			b, err := json.Marshal(msg)
			if err != nil {
				if w.logging {
					w.l.Error("Failed to marshal message", "error", err)
				}
			}

			req := fasthttp.AcquireRequest()
			req.SetRequestURI(w.URL)
			req.SetBody(b)
			req.Header.SetMethod(http.MethodPost)
			req.Header.SetContentType("application/json")

			resp := fasthttp.AcquireResponse()

			err = w.client.Do(req, resp)
			if err != nil {
				if w.logging {
					w.l.Error("Failed to send message", "error", err)
				}
			}
			fasthttp.ReleaseRequest(req)

			code := resp.StatusCode()

			if code != http.StatusOK && code != http.StatusNoContent {
				body := resp.Body()

				if w.logging {
					w.l.Error("Invalid status code", "code", code, "body", string(body))
				}
			}

			fasthttp.ReleaseResponse(resp)
		case <-w.doneChan:
			return
		}
	}
}
