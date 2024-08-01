package models

import (
	"strings"
	"xorm.io/xorm/names"
)

type MdmgTableMapper map[string]bool

var snakeMapper = names.SnakeMapper{}

func (m MdmgTableMapper) Obj2Table(o string) string {
	n := snakeMapper.Obj2Table(o)
	if _, ok := m[n]; ok {
		return "sys_" + n
	}
	return n
}
func (m MdmgTableMapper) Table2Obj(t string) string {
	if strings.HasPrefix(t, "sys_") {
		t = strings.TrimLeft(t, "sys_")
	}
	return snakeMapper.Table2Obj(t)
}
