package project

import (
	"github.com/simplejia/push/model"
)

func (project *Project) Count() (n int, err error) {
	projectModel := model.NewProject()
	n, err = projectModel.Count()
	if err != nil {
		return
	}
	return
}
