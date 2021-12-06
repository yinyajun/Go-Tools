package dict

import "time"

// dict 接口，目前使用cache作为dict的底层实现
type Dict interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{}, expire time.Duration)
	Del(key string) bool
	Stat() string
}
