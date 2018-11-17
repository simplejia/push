package constraint

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (constraint *Constraint) Inc(field string, amount, max int) (updatedConstraint *Constraint, err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"_id": constraint.ID,
	}

	if max > 0 {
		q[field] = bson.M{
			"$lt": max,
		}
	}

	update := bson.M{
		"$inc": bson.M{
			field: amount,
		},
	}

	change := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}

	_, err = c.Find(q).Apply(change, &updatedConstraint)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}

		err = nil
	}

	return
}

func (constraint *Constraint) IncWorkerCounter(amount, max int) (updatedConstraint *Constraint, err error) {
	return constraint.Inc("worker_counter", amount, max)
}

func (constraint *Constraint) IncTargetCounter(amount, max int) (updatedConstraint *Constraint, err error) {
	return constraint.Inc("target_counter", amount, max)
}
