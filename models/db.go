package models

import "xorm.io/xorm"

type Context struct {
	Engine *xorm.Engine
}

func NewContextWithEngine(eng *xorm.Engine) *Context {
	return &Context{Engine: eng}
}

func NewContext(driverName, dataSourceName string) (*Context, error) {
	eng, err := xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	// 对特定的模型进行sys_前缀修改
	eng.SetTableMapper(MdmgTableMapper{
		"column": true,
		"model":  true,
	})
	return NewContextWithEngine(eng), err
}
