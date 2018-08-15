package target

import (
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/mongo"

	mgo "gopkg.in/mgo.v2"
)

type Target api.Target

func (target *Target) Db() (db string) {
	return "push"
}

func (target *Target) Table() (db string) {
	return "target"
}

func (target *Target) GetC() (c *mgo.Collection) {
	db, table := target.Db(), target.Table()
	session := mongo.DBS[db]
	sessionCopy := session.Copy()
	c = sessionCopy.DB(db).C(table)
	return
}

func (target *Target) Regular() (ok bool) {
	if target == nil {
		return
	}

	ok = true
	return
}
