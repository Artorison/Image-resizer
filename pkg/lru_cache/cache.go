package lrucache

import (
	"sync"
)

type Cache interface {
	Set(key string, value any) bool
	Get(key string) (any, bool)
	Clear()
}

type KeyValItem struct {
	key   string
	value any
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[string]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[string]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key string, value any) bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	cacheItem := KeyValItem{key: key, value: value}

	if item, isInCache := lc.items[key]; isInCache {
		item.Value = cacheItem
		lc.queue.MoveToFront(item)
		return isInCache
	}

	item := lc.queue.PushFront(cacheItem)
	lc.items[key] = item

	if lc.queue.Len() > lc.capacity {
		deliteItem := lc.queue.Back()

		rmItem := deliteItem.Value.(KeyValItem)
		delete(lc.items, rmItem.key)

		lc.queue.Remove(deliteItem)
	}
	return false
}

func (lc *lruCache) Get(key string) (any, bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	item, isInCache := lc.items[key]
	if !isInCache {
		return nil, false
	}
	lc.queue.MoveToFront(item)
	return item.Value.(KeyValItem).value, isInCache
}

func (lc *lruCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.items = make(map[string]*ListItem, lc.capacity)
	lc.queue = NewList()
}
