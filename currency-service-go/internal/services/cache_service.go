package services

import (
	"container/list"
	"sync"
	"time"
)

// CacheService é o equivalente Go do CacheService Java, que por sua vez
// envolve o CacheManager do Spring configurado com Caffeine
// (spec: maximumSize=500,expireAfterWrite=15s no application.yaml).
//
// Como Go não tem um equivalente direto ao Spring Cache abstraction,
// implementamos aqui um cache em memória simples com:
//   - TTL por entrada (equivalente a expireAfterWrite)
//   - tamanho máximo com eviction LRU (equivalente a maximumSize)
//
// O cache é organizado por "cacheName" (ex.: "ConvertedValue"), assim como
// o Java permite múltiplos caches nomeados através do mesmo CacheManager.
type CacheService struct {
	mu      sync.Mutex
	maxSize int
	ttl     time.Duration
	caches  map[string]*namedCache
}

type namedCache struct {
	items map[string]*list.Element
	order *list.List // frente = mais recentemente usado
}

type cacheEntry struct {
	key       string
	value     float64
	expiresAt time.Time
}

// NewCacheService cria um CacheService com o tamanho máximo e TTL informados.
// Os valores default (500, 15s) replicam o spec do Caffeine no application.yaml.
func NewCacheService(maxSize int, ttl time.Duration) *CacheService {
	return &CacheService{
		maxSize: maxSize,
		ttl:     ttl,
		caches:  make(map[string]*namedCache),
	}
}

// Get é o equivalente a CacheService.get(cacheName, key) do Java.
// Retorna (0, false) quando a chave não existe ou expirou —
// equivalente a wrapper == null no Java.
func (c *CacheService) Get(cacheName, key string) (float64, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cache, ok := c.caches[cacheName]
	if !ok {
		return 0, false
	}

	elem, ok := cache.items[key]
	if !ok {
		return 0, false
	}

	entry := elem.Value.(*cacheEntry)
	if time.Now().After(entry.expiresAt) {
		cache.order.Remove(elem)
		delete(cache.items, key)
		return 0, false
	}

	cache.order.MoveToFront(elem)
	return entry.value, true
}

// Put é o equivalente a CacheService.put(cacheName, key, value) do Java.
func (c *CacheService) Put(cacheName, key string, value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cache, ok := c.caches[cacheName]
	if !ok {
		cache = &namedCache{
			items: make(map[string]*list.Element),
			order: list.New(),
		}
		c.caches[cacheName] = cache
	}

	if elem, exists := cache.items[key]; exists {
		entry := elem.Value.(*cacheEntry)
		entry.value = value
		entry.expiresAt = time.Now().Add(c.ttl)
		cache.order.MoveToFront(elem)
		return
	}

	entry := &cacheEntry{key: key, value: value, expiresAt: time.Now().Add(c.ttl)}
	elem := cache.order.PushFront(entry)
	cache.items[key] = elem

	if cache.order.Len() > c.maxSize {
		oldest := cache.order.Back()
		if oldest != nil {
			cache.order.Remove(oldest)
			delete(cache.items, oldest.Value.(*cacheEntry).key)
		}
	}
}
