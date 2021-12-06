package dict

import (
	"encoding/json"
	"hash/crc32"
	"sync"
	"time"

	"github.com/robfig/cron"
)

const (
	DefaultCacheSlotNum  = 256
	CacheKeyNotExist     = -2
	CacheKeyNotSetExpire = -1
)

type cacheNode struct {
	// 最近修改时间
	updateTime time.Time
	// 消息过期时间
	expire time.Duration
	// 消息内容
	val interface{}
}

func (n *cacheNode) isExpired() bool {
	if n.expire <= 0 {
		return false
	}
	return time.Now().Sub(n.updateTime) > n.expire
}

type CacheSlot struct {
	sync.RWMutex
	nodes map[string]*cacheNode
}

type Cache struct {
	cronCycle time.Duration
	slotNum   int
	slots     []*CacheSlot
}

func (c *Cache) getSlot(key string) *CacheSlot {
	sum := crc32.ChecksumIEEE([]byte(key))
	return c.slots[sum%uint32(len(c.slots))]
}

// expire 如果小于等于0，那么永不超时
// set 的对象最好给指针类型的，避免没有必要的copy
func (c *Cache) Set(key string, val interface{}, expire time.Duration) {
	slot := c.getSlot(key)

	node := &cacheNode{
		updateTime: time.Now(),
		val:        val,
		expire:     expire,
	}

	slot.Lock()
	defer slot.Unlock()
	slot.nodes[key] = node
}

func (c *Cache) Del(key string) bool {
	slot := c.getSlot(key)

	slot.Lock()
	defer slot.Unlock()
	if _, ok := slot.nodes[key]; ok {
		delete(slot.nodes, key)
		return true
	}

	return false
}

func (c *Cache) Get(key string) (interface{}, bool) {
	slot := c.getSlot(key)
	slot.RLock()
	defer slot.RUnlock()
	if node, ok := slot.nodes[key]; ok {
		if !node.isExpired() {
			return node.val, true
		}
	}
	return nil, false
}

func (c *Cache) TTL(key string) time.Duration {
	slot := c.getSlot(key)
	slot.RLock()
	defer slot.RUnlock()
	if n, ok := slot.nodes[key]; !ok {
		return CacheKeyNotExist
	} else {
		if n.expire <= 0 {
			return CacheKeyNotSetExpire
		} else {
			ttl := n.expire - time.Now().Sub(n.updateTime)
			if ttl > 0 {
				return ttl
			} else {
				return 0
			}
		}
	}
}

func (c *Cache) Exist(key string) bool {
	slot := c.getSlot(key)

	slot.RLock()
	defer slot.RUnlock()
	node, ok := slot.nodes[key]
	return ok && !node.isExpired()
}

func (c *Cache) CleanExpireNodeCron() {
	for _, slot := range c.slots {
		slot.Lock()
		for k, node := range slot.nodes {
			if node.isExpired() {
				delete(slot.nodes, k)
			}
		}
		slot.Unlock()
	}
}

func (c *Cache) Stat() string {
	stat := make(map[int]int, c.slotNum)

	for i, slot := range c.slots {
		stat[i] = len(slot.nodes)
	}

	s, _ := json.Marshal(stat)
	return string(s)
}

func NewCache(cronCycle time.Duration, slotNum int) *Cache {
	if slotNum <= 0 {
		slotNum = DefaultCacheSlotNum
	}
	slots := make([]*CacheSlot, slotNum)
	for i := 0; i < slotNum; i++ {
		slots[i] = &CacheSlot{
			nodes: make(map[string]*cacheNode),
		}
	}
	cache := &Cache{
		cronCycle: cronCycle,
		slotNum:   slotNum,
		slots:     slots,
	}

	// 清理过期数据cron
	c := cron.New()
	c.Schedule(cron.Every(cronCycle), cron.FuncJob(cache.CleanExpireNodeCron))
	c.Start()

	return cache
}
