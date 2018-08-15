package project

import (
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
)

func (project *Project) Get(projectID int64) (projectAPI *api.Project, err error) {
	projectModel := model.NewProject()
	projectModel.ID = projectID
	projectModelRet, err := projectModel.Get()
	if err != nil {
		return
	}

	if projectModelRet == nil {
		return
	}

	projectAPI = (*api.Project)(projectModelRet)

	return
}
