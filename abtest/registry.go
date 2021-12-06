/*
* @Author: Yajun
* @Date:   2021/12/5 15:33
 */

package abtest

import (
	"sync/atomic"
	"unsafe"
)

type Registry struct {
	Name string
	Dict map[string]*unsafe.Pointer
}

func (r *Registry) Lookup(name string) (*Domain, bool) {
	if addr, ok := r.Dict[name]; ok {
		return CurrentDomain(addr), ok
	}
	return nil, false
}

func (r *Registry) Inject(s Storage) error {
	for k, v := range r.Dict {
		if err := s.RegisterFunc(k, UpdateFuncFactory(v)); err != nil {
			return err
		}
	}
	return nil
}

func UpdateFuncFactory(addr *unsafe.Pointer) UpdateFunc {
	return func(key string, data []byte) error {
		domain, err := Parse(key, data)
		if err != nil {
			return err
		}
		atomic.StorePointer(addr, unsafe.Pointer(domain))
		return nil
	}
}

func CurrentDomain(addr *unsafe.Pointer) *Domain {
	return (*Domain)(atomic.LoadPointer(addr))
}
