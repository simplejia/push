package constraint

import (
	"github.com/simplejia/push/model"
)

func (constraint *Constraint) Proceed(id int64) (exist bool, err error) {
	constraintModel := model.NewConstraint()
	constraintModel.ID = id
	exist, err = constraintModel.Proceed()
	if err != nil {
		return
	}

	return
}
