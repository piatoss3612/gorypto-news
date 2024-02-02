package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"sync"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type Cache interface {
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttls ...time.Duration) error
}

type InMemoryCache struct {
	cache map[string][]byte // [key]value

	mu sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[string][]byte),
		mu:    sync.RWMutex{},
	}
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)

	return nil
}

func (c *InMemoryCache) Exists(ctx context.Context, key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.cache[key]

	return ok
}

func (c *InMemoryCache) Get(ctx context.Context, key string, value interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.cache[key]
	if !ok {
		return ErrCacheMiss
	}

	return gob.NewDecoder(bytes.NewBuffer(v)).Decode(value)
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, ttls ...time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	buf := bytes.NewBuffer(nil)

	err := gob.NewEncoder(buf).Encode(value)
	if err != nil {
		return err
	}

	c.cache[key] = buf.Bytes()

	ttl := 1 * time.Hour
	if len(ttls) > 0 {
		ttl = ttls[0]
	}

	if ttl == 0 {
		ttl = 1 * time.Hour
	}

	time.AfterFunc(ttl, func() {
		c.Delete(ctx, key)
	})

	return nil
}
