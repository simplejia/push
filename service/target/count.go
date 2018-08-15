package target

import (
	"github.com/simplejia/push/model"
)

func (target *Target) Count() (n int, err error) {
	targetModel := model.NewTarget()
	n, err = targetModel.Count()
	if err != nil {
		return
	}

	return
}

func (target *Target) CountByCid(cid int64) (n int, err error) {
	targetModel := model.NewTarget()
	n, err = targetModel.CountByCid(cid)
	if err != nil {
		return
	}

	return
}
