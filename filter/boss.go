package filter

import (
	"net/http"
	"runtime/debug"
	"time"

	"lib"

	"github.com/simplejia/clog"
)

// Boss 后置过滤器，用于数据上报，比如调用延时，出错等
func Boss(w http.ResponseWriter, r *http.Request, m map[string]interface{}) bool {
	err := m["__E__"]
	path := m["__P__"]
	c := m["__C__"].(lib.IBase)
	bt := m["__T__"].(time.Time)

	if err != nil {
		clog.Error("Boss() path: %v, body: %s, err: %v, stack: %s", path, c.ReadBody(r), err, debug.Stack())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		clog.Info("Boss() path: %v, body: %s, elapse: %s", path, c.ReadBody(r), time.Since(bt))
	}
	return true
}
