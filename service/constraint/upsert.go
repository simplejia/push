package constraint

import (
	"github.com/simplejia/push/api"
	constraint_model "github.com/simplejia/push/model/constraint"
)

func (constraint *Constraint) Upsert(constraintAPI *api.Constraint) (err error) {
	if constraintAPI == nil {
		return
	}

	constraintModel := (*constraint_model.Constraint)(constraintAPI)
	err = constraintModel.Upsert()
	if err != nil {
		return
	}

	return
}
