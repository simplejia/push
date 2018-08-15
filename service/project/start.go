package project

import (
	"github.com/simplejia/push/model"
)

func (project *Project) Start(id int64) (exist bool, err error) {
	projectModel := model.NewProject()
	projectModel.ID = id
	exist, err = projectModel.Start()
	if err != nil {
		return
	}

	return
}
