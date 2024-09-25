package base

import "xorm.io/xorm"

type IcInitTable interface {
	InitTable(engine *xorm.Engine) error
}
