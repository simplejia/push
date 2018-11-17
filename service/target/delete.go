package target

import (
	"github.com/globalsign/mgo/bson"
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
