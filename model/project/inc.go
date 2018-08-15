package project

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (project *Project) Inc(field string, amount, max int) (updatedProject *Project, err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"_id": project.ID,
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

	_, err = c.Find(q).Apply(change, &updatedProject)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}

		err = nil
	}

	return
}

func (project *Project) IncWorkerCounter(amount, max int) (updatedProject *Project, err error) {
	return project.Inc("worker_counter", amount, max)
}

func (project *Project) IncTargetCounter(amount, max int) (updatedProject *Project, err error) {
	return project.Inc("target_counter", amount, max)
}
