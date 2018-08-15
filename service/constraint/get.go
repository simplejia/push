package constraint

import (
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
)

func (constraint *Constraint) Get(constraintID int64) (constraintRet *api.Constraint, err error) {
	constraintModel := model.NewConstraint()
	constraintModel.ID = constraintID
	constraintModelRet, err := constraintModel.Get()
	if err != nil {
		return
	}

	if constraintModelRet == nil {
		return
	}

	constraintRet = (*api.Constraint)(constraintModelRet)
	return
}
