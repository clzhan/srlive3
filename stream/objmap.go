package stream


import (
"time"

"github.com/clzhan/srlive3/log"
cmap "github.com/streamrail/concurrent-map"
)

var (
	objects = cmap.New()

	//log      *util.FileLogger
)

func AddObject(obj *StreamObject) {
	objects.Set(obj.name, obj)
}

func RemoveObject(name string) {
	objects.Remove(name)
}

func FindObject(name string) (*StreamObject, bool) {
	if v, found := objects.Get(name); found {
		return v.(*StreamObject), true
	}
	return nil, false
}
func Timer() {

	for {
		select {
		case <-time.After(10 * time.Second):
			objects.IterCb(func(key string, v interface{}) {
				log.Info("object key :", key, " value :", v)
			})
		}
	}
}

