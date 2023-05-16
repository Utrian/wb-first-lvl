package cache

import (
	"time"
	"wb-first-lvl/internal/database/queries"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	repo  queries.OrderRepo
	cache *cache.Cache
}

func CreateCache(rep queries.OrderRepo) *Cache {
	return &Cache{
		repo:  rep,
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

// func (c *Cache) ItinCache() error {
// 	ords, err :=
// }
