package project

import (
	"gopkg.in/mgo.v2"
)

func (project *Project) Get() (projectRet *Project, err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	err = c.FindId(project.ID).One(&projectRet)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}

		err = nil
		return
	}

	return
}
