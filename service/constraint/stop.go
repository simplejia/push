package constraint

import (
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
)

func (constraint *Constraint) Stop(id int64) (exist bool, err error) {
	curState := api.StateStarted

	constraintModel := model.NewConstraint()
	constraintModel.ID = id
	exist, err = constraintModel.UpdateState(curState, api.StateStopped)
	if err != nil {
		return
	}

	return
}
