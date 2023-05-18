package cache

import (
	"sync"
	"time"
	"wb-first-lvl/internal/models"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	ords              map[string]Ord
}

type Ord struct {
	Value      models.Order
	Created    time.Time
	Expiration int64
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	ords := make(map[string]Ord)

	cache := Cache{
		ords:              ords,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

func (c *Cache) Set(ord_uid string, order models.Order, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.ords[ord_uid] = Ord{
		Value:      order,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *Cache) Get(ord_uid string) (models.Order, bool) {
	c.RLock()
	defer c.RUnlock()

	ord, found := c.ords[ord_uid]

	if !found {
		return models.Order{}, false
	}

	if ord.Expiration > 0 {
		if time.Now().UnixNano() > ord.Expiration {
			return models.Order{}, false
		}
	}

	return ord.Value, true
}

func (c *Cache) RestoreCache(ords *[]models.Order) {
	for _, ord := range *ords {
		c.Set(ord.OrderUID, ord, 0)
	}
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)

		if c.ords == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearOrds(keys)
		}
	}
}

func (c *Cache) expiredKeys() (keys []string) {
	c.RLock()
	defer c.RUnlock()

	for k, o := range c.ords {
		if time.Now().UnixNano() > o.Expiration && o.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (c *Cache) clearOrds(keys []string) {
	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.ords, k)
	}
}
