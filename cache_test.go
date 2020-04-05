package objcache

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type simpleBytesCache struct {
	data []byte
}

func (sb simpleBytesCache) Read(ctx context.Context, key string) ([]byte, error) {
	return sb.data, nil
}

func (sb simpleBytesCache) GetMap(ctx context.Context, keys ...string) (map[string][]byte, error) {
	m := map[string][]byte{}
	for _, key := range keys {
		m[key] = sb.data
	}
	return m, nil
}

func TestObjCache(t *testing.T) {
	type typ struct {
		A int
	}

	value := &typ{1}

	bytes, _ := json.Marshal(value)

	sb := simpleBytesCache{bytes}

	objCache, err := NewObjCache(
		context.Background(),
		(*typ)(nil),
		WithTTL(1*time.Second),
		WithSize(100),
		WithLoadObjectMap(func(ctx context.Context, keys ...string) (i interface{}, e error) {
			return nil, nil
		}),
		WithMarshal(json.Marshal, json.Unmarshal),
		WithManagedBytesCacheFactory(func(ctx context.Context, size int, ttl time.Duration, loadBytesMap FuncLoadBytesMap) ManagedBytesCache {
			return sb
		}),
	)

	y, err := objCache.GetMap(context.Background(), "abc")
	fmt.Println(err)
	fmt.Printf("%#v\n", y.(map[string]*typ)["abc"])
}

func BenchmarkObjCache(b *testing.B) {
	type typ struct {
		A int
	}

	value := &typ{1}

	bytes, _ := json.Marshal(value)

	sb := simpleBytesCache{bytes}
	_ = sb
}
