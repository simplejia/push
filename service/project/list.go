package project

import (
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
)

func (project *Project) List(offset, num int) (projects []*api.Project, err error) {
	projectModel := model.NewProject()
	projectsModel, err := projectModel.List(offset, num)
	if err != nil {
		return
	}

	for _, projectModel := range projectsModel {
		projects = append(projects, (*api.Project)(projectModel))
	}

	return
}
