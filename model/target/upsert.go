package target

import (
	"github.com/globalsign/mgo/bson"
)

func (target *Target) Upsert() (exist bool, err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	sel := bson.M{
		"cid": target.CID,
		"mid": target.MID,
	}

	up := bson.M{
		"$setOnInsert": bson.M{
			"_id": bson.NewObjectId(),
		},
	}

	info, err := c.Upsert(sel, up)
	if err != nil {
		return
	}

	if info.Matched == 1 {
		exist = true
	}

	return
}
