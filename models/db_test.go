package models

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
	"testing"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var ctx *Context
var dbName string

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	driverName := "mysql"
	if driverName == "mysql" {
		ctx, _ = NewContext("mysql", "root:123456@tcp(127.0.0.1:3306)/mdmg?charset=utf8")
	} else {
		dbName = generateTestDbName()
		removeExistTestDb(dbName)
		ctx, _ = NewContext("sqlite3", dbName)
	}

	ctx.Engine.ShowSQL(true)
	ctx.Engine.Logger().SetLevel(log.LOG_DEBUG)
	// ctx, _ = NewContext("sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
}

func removeExistTestDb(testDbName string) {
	if _, err := os.Stat(testDbName); err == nil {
		_ = os.Remove(testDbName)
		fmt.Println(testDbName + " is exist, remove it")
	}
}

func generateTestDbName() string {
	exe, _ := os.Executable()
	exe = strings.ReplaceAll(exe, string(os.PathSeparator), "_")
	testDbName := fmt.Sprintf("%x.db", md5.Sum([]byte(exe)))
	fmt.Println("testDbName:", testDbName)
	return testDbName
}

func teardown() {
	_ = ctx.Engine.Close()
	// removeExistTestDb(dbName)
}

func TestNewContext(t *testing.T) {
	_ = ctx.Engine.CreateTables(Column{}, Model{})
}

func TestNewContext2(t *testing.T) {
	// ctx.Engine.CreateTables(Column{}, Model{})
}

func TestDropAllTestTables(t *testing.T) {
	err := ctx.Engine.DropTables(Column{}, Model{})
	if err != nil {
		t.Error(err)
	}
	metas, err := ctx.Engine.DBMetas()
	if err != nil || len(metas) > 0 {
		t.Error(err)
	}

}

func TestColumnJsonMarshal(t *testing.T) {
	col := &Column{}
	jd, err := json.Marshal(col)
	if err != nil {
		t.Error(err)
	}
	// schemas.Table{}
	tab, err := ctx.Engine.TableInfo(col)
	if err != nil {
		t.Error(err)
	}
	err = ctx.Engine.CreateTables(*col)
	if err != nil {
		t.Error(err)
	}
	for _, c := range tab.Columns() {
		col.ConvertFrom(c)
		col.ID = 0
		//jd, err = json.Marshal(col)
		//fmt.Println(string(jd))
		//jd, err = json.Marshal(col.SQLType)
		//fmt.Println("sql_type:", string(jd))
		affected, err := ctx.Engine.Insert(col)
		if err != nil {
			t.Error(err)
		}
		t.Logf("return affected num %v\n", affected)
	}

	fmt.Printf("%v", tab)
	fmt.Println(string(jd))
}

func TestDumpMysql(t *testing.T) {
	eng, err := xorm.NewEngine("mysql", "devuser:devuser.COM2019@tcp(devmysql01.wiz.top:6033)/wiz_hr2?charset=utf8")
	if err != nil {
		t.Error(err)
	}
	err = eng.DumpAllToFile("./dump.sql")
	if err != nil {
		t.Error(err)
	}
}

func TestQueryAnyToJson(t *testing.T) {
	eng, err := xorm.NewEngine("mysql", "devuser:devuser.COM2019@tcp(devmysql01.wiz.top:6033)/wiz_hr2?charset=utf8")
	if err != nil {
		t.Error(err)
	}
	sql := "select * from hr_att_duty_work_overtime"
	//ret, err := eng.QueryInterface(sql)
	//if err != nil {
	//	t.Error(err)
	//}
	//jRet, err := json.MarshalIndent(ret, "", " ")
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(string(jRet))

	ret2, err := eng.Exec(sql)
	if err != nil {
		t.Error(err)
	}
	lastInsertId, _ := ret2.LastInsertId()
	affected, _ := ret2.RowsAffected()
	t.Log(lastInsertId, affected)
}
