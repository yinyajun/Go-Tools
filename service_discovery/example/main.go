package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	g "discovery"
)

var (
	endpoint   string
	replicaNum int
	hash       string
)

const (
	Prefix   = "/abc/"
	User     = ""
	Password = ""
)

func init() {
	flag.StringVar(&endpoint, "etcd", "localhost:2379", "etcd endpoints")
	flag.IntVar(&replicaNum, "num", 1000, "replica num")
	flag.StringVar(&hash, "hash", "murmur3", "murmur3, crc32, fnv")
}

func randomCases(num int) []string {
	ret := []string{}
	for i := 0; i < num; i++ {
		base := rand.Intn(9) * 10000000
		delta := rand.Intn(10000000)
		ret = append(ret, strconv.Itoa(base+delta))
	}
	return ret
}

func showStat(stat map[string]int, errNum int) {
	var allNum int
	for _, v := range stat {
		allNum += v
	}
	fmt.Println("-------------------------------------")
	length := float64(len(stat))
	avg := 1.0 / length
	variance := 0.0
	for k, v := range stat {
		freq := float64(v) / float64(allNum)
		fmt.Println(k, freq)
		variance += (freq - avg) * (freq - avg)
	}
	fmt.Println("err num:", errNum)
	fmt.Println("variance:", variance/length)
}

func main() {
	_ = g.DefaultServiceRegister(endpoint, User, Password, Prefix)
	d := g.DefaultServiceDiscovery(endpoint, User, Password, Prefix, replicaNum, hash)

	for {
		select {
		case <-time.Tick(5 * time.Second):
			stat := map[string]int{}
			errNum := 0
			for _, k := range randomCases(100) {
				node, err := d.Allocate(k)
				if err != nil {
					errNum += 1
				}
				stat[node] += 1
			}
			showStat(stat, errNum)
		}
	}
}
