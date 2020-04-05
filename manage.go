package objcache

import (
	"context"
	"fmt"
	"time"
)

type managedBytesCache struct {
	size         int
	ttl          time.Duration
	loadBytesMap FuncLoadBytesMap
	bytesCache   BytesCache
}

func NewManagedBytesCache(ctx context.Context, size int, ttl time.Duration, loadBytesMap FuncLoadBytesMap, bytesCache BytesCache) ManagedBytesCache {
	return &managedBytesCache{
		ttl:          ttl,
		loadBytesMap: loadBytesMap,
		bytesCache:   bytesCache,
	}
}

func (m *managedBytesCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, errNotImplemented
}

func (m *managedBytesCache) GetMap(ctx context.Context, keys ...string) (map[string][]byte, error) {
	bytesMap := make(map[string][]byte, len(keys))
	var missKeys []string

	for _, key := range keys {
		bytes, err := m.bytesCache.Get(ctx, key)
		if err == nil {
			bytesMap[key] = bytes
		} else {
			missKeys = append(missKeys, key)
		}
	}

	missMap, err := m.loadBytesMap(ctx, missKeys...)
	if err != nil {
		return nil, fmt.Errorf("managedBytesCache: %v", err)
	}

	for key, value := range missMap {
		bytesMap[key] = value
		// TODO: handle error
		_ = m.bytesCache.Set(ctx, key, value, m.ttl)
	}

	return bytesMap, nil
}
