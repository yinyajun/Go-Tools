/*
* @Author: Yajun
* @Date:   2021/11/25 16:54
 */

package main

import (
	"fmt"
	"log"
	"time"

	"dag_flow"
)

type Job struct {
	// 节点的唯一标志
	taskID     uint64
	isFinished bool
}

func NewJob(taskID uint64) Job {
	return Job{
		taskID: taskID,
	}
}

func (j *Job) GetTaskID() uint64 {
	return j.taskID
}

func (j *Job) Exec() {
	time.Sleep(time.Second)
	fmt.Printf("taskID[%d] here is exec module\n", j.taskID)
}

func (j *Job) Complete() {
	time.Sleep(time.Second)
	fmt.Printf("taskID[%d] here is complete module\n", j.taskID)
}

func (j *Job) Hashcode() interface{} {
	return j.taskID
}

func (j *Job) IsFinished() bool {
	return j.isFinished
}

func (j *Job) SetFinished(bo bool) {
	j.isFinished = bo
}

func main() {
	var df dag_flow.DagFlow
	job1 := NewJob(1)
	job2 := NewJob(2)
	job3 := NewJob(3)
	job4 := NewJob(4)
	job5 := NewJob(5)
	job6 := NewJob(6)
	df.Add(&job1)
	df.Add(&job2)
	df.Add(&job3)
	df.Add(&job4)
	df.Add(&job5)
	df.Add(&job6)
	df.Connect(&job1, &job2)
	df.Connect(&job1, &job3)
	df.Connect(&job1, &job4)
	df.Connect(&job2, &job5)
	df.Connect(&job3, &job5)
	df.Connect(&job5, &job6)
	df.Connect(&job4, &job6)

	fmt.Println(df.Validate())

	dot := df.Graph().Dot(&dag_flow.DotOpts{})
	fmt.Printf(string(dot))

	err := df.Run()
	if err != nil {
		log.Fatalf("run failed,err:%v", err)
	}
	log.Printf("run ok")
}
