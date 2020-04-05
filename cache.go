package objcache

import (
	"context"
	"fmt"
	"reflect"
)

type ManagedBytesCache interface {
	GetMap(ctx context.Context, keys ...string) (map[string][]byte, error)
}

type objCache struct {
	typ       reflect.Type
	marshal   FuncMarshal
	unmarshal FuncUnmarshal

	managedBytesCache ManagedBytesCache
}

func NewObjCache(ctx context.Context, proto interface{}, options ...Option) (ObjCache, error) {
	opt := newOptions()
	for _, o := range options {
		o(opt)
	}

	if opt.loadObjectMap == nil {
		return nil, fmt.Errorf("NewObjCache: loadObjectMap not provided")
	}

	var managedBytesCache ManagedBytesCache

	if opt.managedBytesCacheFactory != nil {
		managedBytesCache = opt.managedBytesCacheFactory(ctx, opt.size, opt.ttl, makeLoadBytesMap(opt.loadObjectMap, opt.marshal))
	} else if opt.bytesCache != nil {
		managedBytesCache = NewManagedBytesCache(ctx, opt.size, opt.ttl, makeLoadBytesMap(opt.loadObjectMap, opt.marshal), opt.bytesCache)
	} else {
		return nil, fmt.Errorf("NewObjCache: no bytes cache provided")
	}

	return &objCache{
		typ:       reflect.TypeOf(proto),
		marshal:   opt.marshal,
		unmarshal: opt.unmarshal,

		managedBytesCache: managedBytesCache,
	}, nil
}

func (o *objCache) unmarshalWrap(data []byte) (reflect.Value, error) {
	return unmarshalWrap(data, o.typ, o.unmarshal)
}

func (o *objCache) GetMap(ctx context.Context, keys ...string) (interface{}, error) {
	bytesMap, err := o.managedBytesCache.GetMap(ctx, keys...)
	if err != nil {
		return nil, fmt.Errorf("objCache.GetMap: %v", err)
	}

	objMap := reflect.MakeMapWithSize(reflect.MapOf(reflect.TypeOf(""), o.typ), len(bytesMap))
	for key, bytes := range bytesMap {
		value, err := o.unmarshalWrap(bytes)
		if err != nil {
			return nil, fmt.Errorf("objCache.GetMap: [%s]%v", key, err)
		}
		objMap.SetMapIndex(reflect.ValueOf(key), value)
	}
	return objMap.Interface(), nil
}
