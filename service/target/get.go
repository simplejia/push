package target

import (
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
)

func (target *Target) Get(cid, mid int64) (targetAPI *api.Target, err error) {
	targetModel := model.NewTarget()
	targetModel.CID = cid
	targetModel.MID = mid
	targetRet, err := targetModel.Get()
	if err != nil {
		return
	}

	if targetRet == nil {
		return
	}

	targetAPI = (*api.Target)(targetRet)

	return
}
