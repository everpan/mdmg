package config

import (
	"github.com/everpan/mdmg/pkg/base/store"
	"xorm.io/xorm"
)

var (
	GlobalConfig = NewIConfig()
	engineCache  = store.OneLevelMap[string, *xorm.Engine]{}
)

func FetchEngine(connStr string) *xorm.Engine {
	e, ok := engineCache.Get(connStr)
	if ok {
		return e
	}
	return nil
}

func CacheEngine(connStr string, engine *xorm.Engine) {
	engineCache.Set(connStr, engine)
}

func AcquireEngine(driver, connStr string) (*xorm.Engine, error) {
	e := FetchEngine(connStr)
	if e == nil {
		var err error
		e, err = xorm.NewEngine(driver, connStr)
		if err == nil {
			CacheEngine(connStr, e)
		}
		return e, err
	}
	return e, nil
}
