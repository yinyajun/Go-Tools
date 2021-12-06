package core

import "unsafe"

type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

func (s *Stat) add(objKey, objValue *PepperObject) {
	s.Count += 1
	s.KeySize += getSize(objKey)
	s.ValueSize += getSize(objValue)
}

func (s *Stat) del(objKey, objValue *PepperObject) {
	s.Count -= 1
	s.KeySize -= getSize(objKey)
	s.ValueSize -= getSize(objValue)
}

func getSize(obj *PepperObject) int64 {
	switch o := obj.Ptr.(type) {
	case string:
		return int64(len(o))
	default:
		return int64(unsafe.Sizeof(o))
	}
}
