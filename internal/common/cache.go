package common

import (
	"sync"
	"time"
)

// CacheItem representa un elemento en el cache
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// Cache es un cache simple en memoria
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// NewCache crea una nueva instancia de cache
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// Set almacena un valor en el cache con tiempo de expiración
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get obtiene un valor del cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Verificar si ha expirado
	if time.Now().After(item.ExpiresAt) {
		// Limpiar elemento expirado
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		c.mu.RLock()
		return nil, false
	}

	return item.Data, true
}

// Delete elimina un elemento del cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear limpia todo el cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}

// CleanupExpired limpia elementos expirados
func (c *Cache) CleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.ExpiresAt) {
			delete(c.items, key)
		}
	}
}

// Global cache instance
var GlobalCache = NewCache()

// StartCacheCleanup inicia la limpieza periódica del cache
func StartCacheCleanup() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			GlobalCache.CleanupExpired()
		}
	}()
}
