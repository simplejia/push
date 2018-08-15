package project

import (
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
)

func (project *Project) Stop(id int64) (exist bool, err error) {
	curState := api.StateStarted

	projectModel := model.NewProject()
	projectModel.ID = id
	exist, err = projectModel.UpdateState(curState, api.StateStopped)
	if err != nil {
		return
	}

	return
}
