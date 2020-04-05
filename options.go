package objcache

import (
	"context"
	"encoding/json"
	"time"
)

type (
	FuncMarshal   = func(interface{}) ([]byte, error)
	FuncUnmarshal = func([]byte, interface{}) error

	FuncLoadObject    = func(ctx context.Context, key string) (interface{}, error)
	FuncLoadObjectMap = func(ctx context.Context, keys ...string) (interface{}, error)
	FuncLoadBytesMap  = func(ctx context.Context, keys ...string) (map[string][]byte, error)
)

type BytesCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string)
}

type ManagedBytesCacheFactory = func(ctx context.Context, size int, ttl time.Duration, loadBytesMap FuncLoadBytesMap) ManagedBytesCache

type Options struct {
	size int
	ttl  time.Duration

	loadObjectMap FuncLoadObjectMap

	marshal   FuncMarshal
	unmarshal FuncUnmarshal

	bytesCache               BytesCache
	managedBytesCacheFactory ManagedBytesCacheFactory
}

func newOptions() *Options {
	// TODO: add default bytes cache

	return &Options{
		marshal:   json.Marshal,
		unmarshal: json.Unmarshal,
	}
}

type Option func(*Options)

func WithSize(size int) Option {
	return func(options *Options) {
		options.size = size
	}
}

func WithTTL(ttl time.Duration) Option {
	return func(options *Options) {
		options.ttl = ttl
	}
}

func WithLoadObjectMap(loadObjectMap FuncLoadObjectMap) Option {
	return func(options *Options) {
		options.loadObjectMap = loadObjectMap
	}
}

func WithMarshal(marshal FuncMarshal, unmarshal FuncUnmarshal) Option {
	return func(options *Options) {
		options.marshal = marshal
		options.unmarshal = unmarshal
	}
}

func WithBytesCache(bytesCache BytesCache) Option {
	return func(options *Options) {
		options.bytesCache = bytesCache
	}
}

func WithManagedBytesCacheFactory(factory ManagedBytesCacheFactory) Option {
	return func(options *Options) {
		options.managedBytesCacheFactory = factory
	}
}
