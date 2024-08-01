package utils

import (
	"strings"
	"xorm.io/xorm/schemas"
)

// FilterTables 过滤部分数据库
func FilterTables(tables []*schemas.Table) []*schemas.Table {
	// 1. 如果数据库名以 _开始，视为系统功能表，隐藏
	var result = make([]*schemas.Table, 0, len(tables))
	for _, table := range tables {
		if strings.HasPrefix(table.Name, "_") {
			continue
		}
		result = append(result, table)
	}
	return result
}
