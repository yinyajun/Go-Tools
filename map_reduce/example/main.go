/*
* @Author: Yajun
* @Date:   2021/12/12 20:50
 */

package main

import (
	"fmt"
	"map_reduce"
	"time"
)

func NumGenerator() chan interface{} {
	gen := make(chan interface{})
	go func() {
		for i := 0; i < 20; i++ {
			gen <- i
		}
		close(gen)
	}()
	return gen
}

func reducer3(input chan interface{}, output chan interface{}) {
	res := map[int]int{}
	for c := range input {
		for k, v := range c.(map[int]int) {
			res[k] = v
		}

	}
	output <- res
}

func mapper3(item interface{}, out chan interface{}) {
	fmt.Println("mapper", item)
	res := map[int]int{}
	n := item.(int)
	res[n] = n * 2
	time.Sleep(5 * time.Second)
	out <- res
}

func main() {
	res := map_reduce.MapReduceTask(mapper3, reducer3, NumGenerator(), 3)
	fmt.Println(res)
}
