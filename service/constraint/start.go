package constraint

import (
	"github.com/simplejia/push/model"
)

func (constraint *Constraint) Start(id int64) (exist bool, err error) {
	constraintModel := model.NewConstraint()
	constraintModel.ID = id
	exist, err = constraintModel.Start()
	if err != nil {
		return
	}

	return
}
