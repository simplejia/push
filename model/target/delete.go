package target

import (
	"github.com/globalsign/mgo/bson"
)

func (target *Target) DeletePids() (err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	sel := bson.M{
		"cid": target.CID,
		"mid": target.MID,
	}

	update := bson.M{
		"$unset": bson.M{
			"pids": nil,
		},
	}

	err = c.Update(sel, update)
	if err != nil {
		return
	}

	return
}

func (target *Target) Delete() (err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	err = c.RemoveId(target.ID)
	if err != nil {
		return
	}

	return
}
