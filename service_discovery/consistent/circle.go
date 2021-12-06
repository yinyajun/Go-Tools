// modified "stathat.com/c/consistent"
// 1. use treemap to replace map (n^2klognk -> nklognk)
// 2. support murmur3 hash

package consistent

import (
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"strconv"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/spaolacci/murmur3"
)

type Circle struct {
	sync.RWMutex
	*treemap.Map
	ReplicaNum int
	members    map[string]struct{}
	count      int
	Hash       string
}

func NewCircle() *Circle {
	c := new(Circle)
	c.Map = treemap.NewWith(utils.UInt32Comparator)
	c.ReplicaNum = 20
	c.members = make(map[string]struct{})
	return c
}

func (c *Circle) virtualNode(realNode string, idx int) string {
	//hash := strconv.Itoa(int(murmur3.Sum64([]byte(realNode))))
	return realNode + strconv.Itoa(idx)
}

func (c *Circle) hashKey(key string) uint32 {
	switch c.Hash {
	case "murmur3":
		return c.hashKeyMurmur(key)
	case "crc32":
		return c.hashKeyCRC32(key)
	case "fnv":
		return c.hashKeyFnv(key)
	default:
		return c.hashKeyCRC32(key)
	}
}

func (c *Circle) hashKeyCRC32(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Circle) hashKeyFnv(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (c *Circle) hashKeyMurmur(key string) uint32 {
	h := murmur3.New32()
	h.Write([]byte(key))
	return h.Sum32()
}

func (c *Circle) add(realNode string) {
	for i := 0; i < c.ReplicaNum; i++ {
		vn := c.virtualNode(realNode, i)
		c.Put(c.hashKey(vn), realNode)
	}
	c.members[realNode] = struct{}{}
	c.count++
}

func (c *Circle) Add(realNode string) {
	c.Lock()
	defer c.Unlock()
	c.add(realNode)
}

func (c *Circle) del(realNode string) {
	for i := 0; i < c.ReplicaNum; i++ {
		vn := c.virtualNode(realNode, i)
		c.Remove(c.hashKey(vn))
	}
	delete(c.members, realNode)
	c.count--
}

func (c *Circle) Del(realNode string) {
	c.Lock()
	defer c.Unlock()
	c.del(realNode)
}

func (c *Circle) Members() []string {
	c.RLock()
	defer c.RUnlock()
	var m []string
	for k := range c.members {
		m = append(m, k)
	}
	return m
}

func (c *Circle) Allocate(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	if c.Empty() {
		return "", fmt.Errorf("empty circle")
	}
	key := c.hashKey(name)
	node := c.search(key)
	return node, nil
}

func (c *Circle) search(key uint32) (node string) {
	_, n := c.Ceiling(key)
	if n == nil {
		_, n = c.Min()
	}
	node, _ = n.(string)
	return
}
