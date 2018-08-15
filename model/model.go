/*
Package model 用于模型层定义，所有db及cache对象封装均定义在这里。
只允许在这里添加对外暴露的接口
*/
package model

import (
	"github.com/simplejia/push/model/constraint"
	"github.com/simplejia/push/model/project"
	"github.com/simplejia/push/model/target"
)

// NewConstraint 构造Constraint对象
func NewConstraint() *constraint.Constraint {
	return &constraint.Constraint{}
}

func NewProject() *project.Project {
	return &project.Project{}
}

func NewTarget() *target.Target {
	return &target.Target{}
}
