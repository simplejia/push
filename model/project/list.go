package project

import (
	"github.com/simplejia/push/api"
	"gopkg.in/mgo.v2/bson"
)

func (project *Project) List(offset, limit int) (projects []*Project, err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	err = c.Find(nil).Sort("-_id").Skip(offset).Limit(limit).All(&projects)
	if err != nil {
		return
	}

	return
}

func (project *Project) ListByState(offset, limit int, state api.State) (projects []*Project, err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"state": state,
	}

	err = c.Find(q).Sort("-_id").Skip(offset).Limit(limit).All(&projects)
	if err != nil {
		return
	}

	return
}

func (project *Project) ListAllByState(state api.State) (projects []*Project, err error) {
	return project.ListByState(0, 0, state)
}
