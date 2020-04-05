package objcache

import (
	"context"
	"fmt"
	"reflect"
)

var errNotImplemented = fmt.Errorf("not implemented")

func makeLoadBytesMap(loadObjects FuncLoadObjectMap, marshal FuncMarshal) FuncLoadBytesMap {
	return func(ctx context.Context, keys ...string) (bytes map[string][]byte, e error) {
		m, err := loadObjects(ctx, keys...)
		if err != nil {
			return nil, fmt.Errorf("makeLoadBytesMap: %v", err)
		}

		mValue := reflect.ValueOf(m)
		keyType := mValue.Type().Key()
		if keyType.Kind() != reflect.String {
			return nil, fmt.Errorf("makeLoadBytesMap: key of mapGetter is not string(%v)", keyType)
		}

		bytesMap := make(map[string][]byte, mValue.Len())
		iter := mValue.MapRange()
		for iter.Next() {
			key := iter.Key().String()
			value, err := marshal(iter.Value().Interface())
			if err != nil {
				return nil, fmt.Errorf("makeLoadBytesMap: failed to marshal on key(%s)", key)
			}
			bytesMap[key] = value
		}
		return bytesMap, nil
	}
}

func unmarshalWrap(data []byte, typ reflect.Type, unmarshal FuncUnmarshal) (reflect.Value, error) {
	value := reflect.New(typ)
	err := unmarshal(data, value.Interface())
	return value.Elem(), err
}
