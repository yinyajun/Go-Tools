/*
* @Author: Yajun
* @Date:   2021/12/12 20:48
 */

package map_reduce

import (
	"fmt"
)

type Mapper func(item interface{}, out chan interface{})
type Reducer func(input chan interface{}, output chan interface{})

func MapReduceTask(mapper Mapper, reducer Reducer, source chan interface{}, poolSize int) interface{} {
	reduceInput := make(chan interface{})
	reduceOutput := make(chan interface{})
	workerPool := make(chan chan interface{}, poolSize)

	go reducer(reduceInput, reduceOutput)

	go func() {
		for workerChan := range workerPool {
			reduceInput <- <-workerChan
		}
		close(reduceInput)
	}()

	go func() {
		for item := range source {
			mapOut := make(chan interface{})
			go mapper(item, mapOut)
			fmt.Println("put", item, "mapper result into worker pool")
			workerPool <- mapOut
		}
		close(workerPool)
	}()

	return <-reduceOutput
}
