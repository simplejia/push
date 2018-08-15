package target

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/simplejia/push/model"
)

func (target *Target) Delete(id string) (err error) {
	targetModel := model.NewTarget()
	targetModel.ID = bson.ObjectIdHex(id)
	err = targetModel.Delete()
	if err != nil {
		return
	}

	return
}
