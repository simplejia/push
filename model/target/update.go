package target

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

func (target *Target) SetPidState() (err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	sel := bson.M{
		"cid": target.CID,
		"mid": target.MID,
	}

	set := bson.M{}

	for pid, state := range target.PIDs {
		set[fmt.Sprintf("pids.%s", pid)] = state
	}

	update := bson.M{
		"$set": set,
	}

	err = c.Update(sel, update)
	if err != nil {
		return
	}

	return
}
