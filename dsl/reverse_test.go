package dsl

import (
	"github.com/everpan/mdmg/utils"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"xorm.io/xorm"
)

func TestTableToJson(t *testing.T) {
	t.Skip("manual run")
	orm, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/wiz_hr2?charset=utf8")
	tDB, err := xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if nil != err {
		t.Error(err)
	}
	tables, err := orm.DBMetas()
	tables = utils.FilterTables(tables)

	if nil != err {
		t.Error(err)
	}
	// var colMapper = utils.GetMapperByName("snake")

	for _, table := range tables[0:3] {
		meta := ModelMeta{}
		meta.FromSchemaTable(table)
		mJson, err := meta.Json(true)
		fatalErr(t, err)
		t.Log(mJson)
		//for _, column := range table.Columns() {
		//	colName := colMapper.Table2Obj(column.Name)
		//	colName = ""
		//	//t.Logf("column:%s -> %s", column.Name, colName)
		//}
		err = meta.CreateTable(tDB)
		fatalErr(t, err)
	}

}
