package utils

import (
	"github.com/allegro/bigcache"
	"github.com/jinzhu/gorm"
)

var (
	handler handlerCache
)

type handlerCache struct {
	db            *gorm.DB
	cache         *bigcache.BigCache
	isDevelopment bool
}

func NewCache(cache *bigcache.BigCache, db *gorm.DB, env bool) {

	handler = handlerCache{
		cache:         cache,
		db:            db,
		isDevelopment: env,
	}

}

func GetCache() *bigcache.BigCache {
	return handler.cache
}
