/*
* @Author: Yajun
* @Date:   2021/11/25 18:04
 */

package dag_flow

type JobNode interface {
	// 事件处理
	Exec()
	// 事件处理完成后调用
	Complete()
	// 函数唯一编号
	Hashable
	// 是否完成
	IsFinished() bool
	// 设置结果
	SetFinished(bo bool)
}
