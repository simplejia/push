package project

import (
	"github.com/simplejia/push/api"
	project_model "github.com/simplejia/push/model/project"
)

func (project *Project) Upsert(projectAPI *api.Project) (err error) {
	if projectAPI == nil {
		return
	}

	projectModel := (*project_model.Project)(projectAPI)
	err = projectModel.Upsert()
	if err != nil {
		return
	}

	return
}
