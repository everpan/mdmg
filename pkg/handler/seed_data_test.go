package handler

import (
	"github.com/everpan/mdmg/pkg/base/tenant"
	"os"
	"testing"
	"xorm.io/xorm"
)

func chkError(err error) {
	if nil != err {
		panic(err)
	}
}

func CreateSeedDataSqlite3Engine(dbFile string, rmBefore bool) *xorm.Engine {
	if rmBefore {
		_ = os.Remove(dbFile)
	}
	engine, err := xorm.NewEngine("sqlite3", dbFile)
	chkError(err)
	engine.ShowSQL(true)
	return engine
}

// 构建种子数据，用于测试
func TestBuildSeedDataForTest(t *testing.T) {
	dbFile := "seed_data_test.db"
	engine := CreateSeedDataSqlite3Engine(dbFile, true)
	// 构建租户表
	tInst := &tenant.IcTenantInfo{}
	err := tInst.InitTable(engine)
	chkError(err)
	// 构建默认租户1，2

	tenant.SetSysEngine(engine)
	tenant.DefaultInfo.Idx = 0
	tenant.DefaultInfo.Driver = "sqlite3"
	tenant.DefaultInfo.ConnectString = dbFile
	tenant.DefaultInfo.Save()
	tenant.DefaultHostInfo.Idx = 0
	tenant.DefaultHostInfo.Driver = "sqlite3"
	tenant.DefaultHostInfo.ConnectString = dbFile
	tenant.DefaultHostInfo.Save()
}
