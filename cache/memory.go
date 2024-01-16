package cache

import (
	"errors"
	"github.com/asaskevich/EventBus"
	"sync"
	"time"
)

var nilTime = time.Time{}

var (
	ErrAlreadyLocked = errors.New("already locked")
	ErrKeyNotExist   = errors.New("key not exist")
)

type MemoryCache struct {
	mutexLockMap sync.Map
	cacheMap     sync.Map
	evbus        EventBus.Bus
}

type item struct {
	value     interface{}
	expiredAt time.Time
	count     int
}

func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		mutexLockMap: sync.Map{},
		cacheMap:     sync.Map{},
		evbus:        EventBus.New(),
	}

	// gc
	go cache.gc()
	return cache
}

func (c *MemoryCache) gc() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		c.cacheMap.Range(func(key, value interface{}) bool {
			if value.(*item).expiredAt != nilTime && time.Now().After(value.(*item).expiredAt) {
				c.cacheMap.Delete(key)
			}
			return true
		})
	}
}

func (c *MemoryCache) Get(key string) (interface{}, error) {
	if v, ok := c.cacheMap.Load(key); ok {
		if v.(*item).expiredAt != nilTime && time.Now().After(v.(*item).expiredAt) {
			c.cacheMap.Delete(key)
			return nil, ErrKeyNotExist
		}
		v.(*item).count++
		return v.(*item).value, nil
	}
	return nil, ErrKeyNotExist
}

func (c *MemoryCache) Set(key string, value interface{}) error {
	c.cacheMap.Store(key, &item{value: value, expiredAt: nilTime, count: 0})
	return nil
}

func (c *MemoryCache) SetEx(key string, value interface{}, duration time.Duration) error {
	c.cacheMap.Store(key, &item{value: value, expiredAt: time.Now().Add(duration), count: 0})
	return nil
}

func (c *MemoryCache) Delete(key string) error {
	c.cacheMap.Delete(key)
	return nil
}

func (c *MemoryCache) TryLock(key string) error {
	if _, locked := c.mutexLockMap.LoadOrStore(key, true); locked {
		return ErrAlreadyLocked
	}
	return nil
}

func (c *MemoryCache) Unlock(key string) error {
	c.mutexLockMap.Delete(key)
	return nil
}

func (c *MemoryCache) Subscribe(key string, callback interface{}) error {
	return c.evbus.Subscribe(key, callback)
}

func (c *MemoryCache) Publish(key string, data interface{}) {
	c.evbus.Publish(key, data)
}
