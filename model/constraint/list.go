package constraint

import (
	"github.com/globalsign/mgo/bson"
	"github.com/simplejia/push/api"
)

func (constraint *Constraint) List(offset, num int) (constraints []*Constraint, err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	err = c.Find(nil).Sort("-_id").Skip(offset).Limit(num).All(&constraints)
	if err != nil {
		return
	}

	return
}

func (constraint *Constraint) ListByState(offset, num int, state api.State) (constraints []*Constraint, err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"state": state,
	}

	err = c.Find(q).Sort("-_id").Skip(offset).Limit(num).All(&constraints)
	if err != nil {
		return
	}

	return
}

func (constraint *Constraint) ListAllByState(state api.State) (constraints []*Constraint, err error) {
	return constraint.ListByState(0, 0, state)
}
