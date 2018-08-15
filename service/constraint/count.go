package constraint

import (
	"github.com/simplejia/push/model"
)

func (constraint *Constraint) Count() (n int, err error) {
	constraintModel := model.NewConstraint()
	n, err = constraintModel.Count()
	if err != nil {
		return
	}
	return
}
