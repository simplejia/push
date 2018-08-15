package constraint

import (
	"gopkg.in/mgo.v2"
)

func (constraint *Constraint) Get() (constraintRet *Constraint, err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	err = c.FindId(constraint.ID).One(&constraintRet)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}

		err = nil
	}

	return
}
