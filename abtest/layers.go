/*
* @Author: Yajun
* @Date:   2021/12/3 17:07
 */

package abtest

import (
	"hash/crc32"
	"sort"

	"github.com/spaolacci/murmur3"
)

// Layer 流量层
type Layer struct {
	Name      string     `json:"name" validate:"required,min=2" label:"[层名称]"`
	Traffics  []*Traffic `json:"traffics" validate:"required,unique=Name,gte=1,dive" label:"[流量组]"`
	WhiteList map[ID]int `json:"white_list,omitempty" label:"[白名单]"`
	domain    *Domain
	seed      uint32
	cum       []int
}

func (l *Layer) adjustBuckets() {
	l.cum = make([]int, len(l.Traffics))
	l.cum[0] = int(l.Traffics[0].Share)
	for i := 1; i < len(l.Traffics); i++ {
		l.cum[i] = int(l.Traffics[i].Share) + l.cum[i-1]
	}
}

func (l *Layer) init(d *Domain) {
	var dName string
	if d != nil {
		l.domain = d
		dName = d.name
	}
	l.adjustBuckets()
	if l.seed == 0 {
		l.seed = crc32.ChecksumIEEE([]byte(dName + l.Name))
	}
}

func (l *Layer) shunt(id ID) int {
	m := murmur3.Sum32WithSeed([]byte(id), l.seed)
	percent := float64(m) / MaxUint32
	return int(percent*100) + 1
}

func (l *Layer) getExperiment(idx int, opts ...ReasonOpt) (*Traffic, Reason) {
	var exp *Traffic
	if idx >= 0 && idx < len(l.Traffics) {
		exp = l.Traffics[idx]
	} else {
		exp = l.Traffics[0]
	}
	reason := make(Reason, 0)
	for r := range opts {
		reason = opts[r](reason)
	}
	return exp, reason
}

func (l *Layer) index(bucket int) ReasonOpt {
	if bucket >= 0 && bucket < len(l.Traffics) {
		return validIndex(bucket)
	}
	return inValidIndex
}

func (l *Layer) Allocate(id ID) (*Traffic, Reason) {
	if idx, ok := l.WhiteList[id]; ok {
		return l.getExperiment(idx, whiteList, l.index(idx))
	}
	bucket := l.shunt(id)
	idx := sort.SearchInts(l.cum, bucket)
	return l.getExperiment(bucket, shunt(bucket), l.index(idx))
}
