package objcache

import (
	"context"
)

type ObjCache interface {
	GetMap(ctx context.Context, keys ...string) (interface{}, error)
}
