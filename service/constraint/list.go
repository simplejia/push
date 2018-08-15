package constraint

import (
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
)

func (constraint *Constraint) List(offset, num int) (constraints []*api.Constraint, err error) {
	constraintModel := model.NewConstraint()
	constraintsModel, err := constraintModel.List(offset, num)
	if err != nil {
		return
	}

	for _, constraintModel := range constraintsModel {
		constraints = append(constraints, (*api.Constraint)(constraintModel))
	}

	return
}
