package target

import (
	"github.com/globalsign/mgo/bson"
)

func (target *Target) Count() (n int, err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	n, err = c.Find(nil).Count()
	return
}

func (target *Target) CountByCid(cid int64) (n int, err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"cid": cid,
	}
	n, err = c.Find(q).Count()
	return
}
