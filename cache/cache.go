package cache

import (
	"fmt"
	"runtime"
)

type Cache interface {
	Set(key string, val interface{})
	Get(key string) interface{}
	Del(key string)
	DelOldest()
	Len() int
}

type Value interface {
	Len() int
}

func CalcLen(value interface{}) int {
	switch v := value.(type) {
	case string:
		if runtime.GOARCH == "amd64" {
			return 16 + len(v)
		} else {
			return 8 + len(v)
		}
	case bool, int8, uint8:
		return 1
	case int16, uint16:
		return 2
	case int32, uint32, float32:
		return 4
	case int64, uint64, float64:
		return 8
	case int, uint:
		if runtime.GOARCH == "amd64" {
			return 8
		} else {
			return 4
		}
	case complex64:
		return 8
	case complex128:
		return 16
	case Value:
		return v.Len()
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}
