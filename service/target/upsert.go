package target

import (
	"github.com/simplejia/push/api"
	target_model "github.com/simplejia/push/model/target"
)

func (target *Target) Upsert(targetAPI *api.Target) (exist bool, err error) {
	if targetAPI == nil {
		return
	}

	targetModel := (*target_model.Target)(targetAPI)
	exist, err = targetModel.Upsert()
	if err != nil {
		return
	}

	return
}
