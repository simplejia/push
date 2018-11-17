package constraint

import (
	"github.com/globalsign/mgo"
)

func (constraint *Constraint) Delete() (err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	err = c.RemoveId(constraint.ID)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}

		err = nil
		return
	}

	return
}
