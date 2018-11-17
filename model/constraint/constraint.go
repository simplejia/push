/*
Package constraint offers connections to DB
*/
package constraint

import (
	"time"

	"github.com/simplejia/push/api"
	"github.com/simplejia/push/mongo"

	mgo "github.com/globalsign/mgo"
)

type Constraint api.Constraint

//NewConstraint creates a new constraint pointer
func NewConstraint() *Constraint {
	return &Constraint{}
}

//Db returns constraint db name
func (constraint *Constraint) Db() (db string) {
	return "push"
}

//Table returns constraint table name
func (constraint *Constraint) Table() (table string) {
	return "constraint"
}

func (constraint *Constraint) GetC() (c *mgo.Collection) {
	db, table := constraint.Db(), constraint.Table()
	session := mongo.DBS[db]
	sessionCopy := session.Copy()
	c = sessionCopy.DB(db).C(table)
	return
}

func (constraint *Constraint) Regular() (ok bool) {
	if constraint == nil {
		return
	}

	ok = true
	return
}

func (constraint *Constraint) CheckAllWorkers() (hasAlive, hasDead bool) {
	if len(constraint.Timestamps) == 0 {
		hasDead = true
		return
	}

	now := time.Now().Unix()
	for _, timestamp := range constraint.Timestamps {
		if now-timestamp < int64(api.KeepaliveTimeDeadline/time.Second) {
			hasAlive = true // at least one is alive
		} else {
			hasDead = true // at least one is dead
		}
	}

	return
}
