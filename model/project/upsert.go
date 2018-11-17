package project

import (
	mgo "github.com/globalsign/mgo"
)

func (project *Project) Upsert() (err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	if project.ID <= 0 {
		var _project *Project
		err = c.Find(nil).Sort("-_id").Limit(1).One(&_project)
		if err != nil {
			if err != mgo.ErrNotFound {
				return
			}
			project.ID = 100000
			err = nil
		} else {
			project.ID = _project.ID + 1
		}
	}

	_, err = c.UpsertId(project.ID, project)
	if err != nil {
		return
	}

	return
}
