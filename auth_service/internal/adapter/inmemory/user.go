package inmemory

import (
    "github.com/recktt77/JobFree/internal/model"
    "context"
    "sync"
    "time"
)

type cacheItem struct {
    user      *model.User
    expiresAt time.Time
}

type InMemoryCache struct {
    mu    sync.RWMutex
    store map[string]cacheItem
    ttl   time.Duration
}

func NewInMemoryCache(ttl time.Duration) *InMemoryCache {
    return &InMemoryCache{
        store: make(map[string]cacheItem),
        ttl:   ttl,
    }
}

func (c *InMemoryCache) GetUser(ctx context.Context, id string) (*model.User, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, ok := c.store[id]
    if !ok || time.Now().After(item.expiresAt) {
        return nil, context.Canceled
    }

    return item.user, nil
}

func (c *InMemoryCache) SetUser(ctx context.Context, user *model.User) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.store[user.ID.Hex()] = cacheItem{
        user:      user,
        expiresAt: time.Now().Add(c.ttl),
    }

    return nil
}
