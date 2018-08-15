package project

import (
	"time"

	"github.com/simplejia/push/api"
	"github.com/simplejia/push/mongo"

	mgo "gopkg.in/mgo.v2"
)

type Project api.Project

func NewProject() *Project {
	return &Project{}
}

func (project *Project) Db() (db string) {
	return "push"
}

func (project *Project) Table() (table string) {
	return "project"
}

func (project *Project) GetC() (c *mgo.Collection) {
	db, table := project.Db(), project.Table()
	session := mongo.DBS[db]
	sessionCopy := session.Copy()
	c = sessionCopy.DB(db).C(table)
	return
}

func (project *Project) Regular() (ok bool) {
	if project == nil {
		return
	}

	ok = true
	return
}

func (project *Project) CheckAllWorkers() (hasAlive, hasDead bool) {
	if len(project.Timestamps) == 0 {
		hasDead = true
		return
	}

	now := time.Now().Unix()
	for _, timestamp := range project.Timestamps {
		if now-timestamp < int64(api.KeepaliveTimeDeadline/time.Second) {
			hasAlive = true // at least one is alive
		} else {
			hasDead = true // at least one is dead
		}
	}

	return
}
