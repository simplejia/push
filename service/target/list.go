package target

import (
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
)

func (target *Target) List(offset, num int) (targets []*api.Target, err error) {
	targetModel := model.NewTarget()
	targetsModel, err := targetModel.List(offset, num)
	if err != nil {
		return
	}

	for _, targetModel := range targetsModel {
		targets = append(targets, (*api.Target)(targetModel))
	}

	return
}

func (target *Target) ListByCid(cid int64, offset, num int) (targets []*api.Target, err error) {
	targetModel := model.NewTarget()
	targetsModel, err := targetModel.ListByCid(cid, offset, num)
	if err != nil {
		return
	}

	for _, targetModel := range targetsModel {
		targets = append(targets, (*api.Target)(targetModel))
	}

	return
}
