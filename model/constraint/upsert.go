package constraint

import (
	mgo "gopkg.in/mgo.v2"
)

func (constraint *Constraint) Upsert() (err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	if constraint.ID <= 0 {
		var _constraint *Constraint
		err = c.Find(nil).Sort("-_id").Limit(1).One(&_constraint)
		if err != nil {
			if err != mgo.ErrNotFound {
				return
			}
			constraint.ID = 100000
			err = nil
		} else {
			constraint.ID = _constraint.ID + 1
		}
	}

	_, err = c.UpsertId(constraint.ID, constraint)
	if err != nil {
		return
	}

	return
}
