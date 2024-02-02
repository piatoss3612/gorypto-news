# gorypto news

- 주기적으로 크롤링한 블록체인 관련 뉴스를 인공지능을 통해 요약한 내용을 디스코드 웹훅을 통해 웹훅이 설정된 서버로 전송하는 프로젝트입니다.

## Table of Contents

- [Requirements](#requirements)
- [How to run](#how-to-run)
- [References](#references)

## Requirements

- [Go 1.21.6](https://golang.org/)
- [Docker](https://www.docker.com/)
- [OpenAI](https://openai.com/)
- [Discord](https://discord.com/)
- [Fly.io](https://fly.io/)
- [GitHub Actions](https://docs.github.com/ko/actions)
- [Makefile](https://www.gnu.org/software/make/)

### Packages in use

- [gocolly](https://github.com/gocolly/colly) : web scraping
- [langchaingo](https://github.com/tmc/langchaingo) : natural language processing
- [gocron](https://github.com/go-co-op/gocron) : cron job

## How to run

### Create `.env` file

- `.env.example` 또는 아래 내용을 복사하여 `.env` 파일을 생성합니다.

```bash
WEBHOOK_URL=<YOUR_WEBHOOK_URL>
OPENAI_API_KEY=<YOUR_OPENAI_API_KEY>
OPENAI_MODEL=gpt-3.5-turbo
```

### Run on local

```bash
$ make run
```

or

```bash
$ export $(cat .env | xargs) && go run .
```

### Build and run with Docker

```bash
$ make up
```

or

```bash
$ docker compose up
```

## References

- [TOKENPOST](https://www.tokenpost.kr/)