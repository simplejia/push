package api

import "time"

// 定义 const 常量

// Code 用于定义返回码(detail)
type Code int

// 定义各种返回码(detail)
const (
	CodeProjectNotExist Code = iota + 1
	CodeConstraintNotExist
)

// CodeMap 定义返回码对应的描述
var CodeMap = map[Code]string{
	CodeProjectNotExist:    "推送任务不存在",
	CodeConstraintNotExist: "获取用户任务不存在",
}

const (
	KeepaliveTimeThreshold time.Duration = time.Second
	KeepaliveTimeDeadline  time.Duration = time.Minute
	ErrorThreshold                       = 100
)

// State 定义状态
type State int

const (
	StateNone     State = iota // 0
	StateReady                 // 1
	StateStarted               // 2
	StateStopped               // 3
	StateFinished              // 4
	StateFailed                // 5
	StateReboot                // 6
	StatePause                 // 7
)

type ConditionKind int

const (
	ConditionKindNone     ConditionKind = iota // 0
	ConditionKindConcrete                      // 1
)
