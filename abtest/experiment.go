/*
* @Author: Yajun
* @Date:   2021/12/3 16:25
 */

package abtest

import (
	"fmt"
)

const (
	MaxUint32 = 1<<32 - 1
)

type ID string

// Traffic 流量
type Traffic struct {
	Name  string `json:"name" validate:"required,min=2" label:"[流量实验名称]"`
	Share uint   `json:"share" validate:"required,lte=100" label:"[流量份额]"`
}

// Parameter 实验参数
type Parameter map[string]string

// Reason 分配实验原因
type Reason []string

type ReasonOpt func(Reason) Reason

func shunt(bucket int) ReasonOpt {
	return func(r Reason) Reason {
		return append(r, fmt.Sprintf("bucket:%d", bucket))
	}
}

func whiteList(r Reason) Reason {
	return append(r, "white_list")
}

func validIndex(index int) ReasonOpt {
	return func(r Reason) Reason {
		return append(r, fmt.Sprintf("index:%d", index))
	}
}

func inValidIndex(r Reason) Reason {
	return append(r, fmt.Sprintf("invalid_index"))
}
