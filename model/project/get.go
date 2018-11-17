package project

import (
	"github.com/globalsign/mgo"
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
