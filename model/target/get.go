package target

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (target *Target) Get() (targetRet *Target, err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"cid": target.CID,
		"mid": target.MID,
	}

	err = c.Find(q).One(&targetRet)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}
		err = nil
	}

	return
}
