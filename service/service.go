/*
Package service 用于定义服务层代码。
只允许在这里添加对外暴露的接口
*/
package service

import (
	"github.com/simplejia/push/service/constraint"
	"github.com/simplejia/push/service/project"
	"github.com/simplejia/push/service/target"
)

func NewConstraint() *constraint.Constraint {
	return &constraint.Constraint{}
}

func NewProject() *project.Project {
	return &project.Project{}
}

func NewTarget() *target.Target {
	return &target.Target{}
}
