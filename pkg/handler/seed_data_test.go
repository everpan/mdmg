package handler

import (
	"github.com/everpan/mdmg/pkg/base/entity"
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
	// entity
	eInst := entity.NewContext(engine, 1)
	eInst.InitTable(engine)
	// entity data
	user := &entity.IcEntityClass{ClassName: "user", ClassDesc: "用户信息", PkColumn: "user_id", TenantId: 1}
	eClass, _ := eInst.RegisterEntityClass(user)
	uAcc := struct {
		UserId       uint32 `xorm:"pk autoincr not null"`
		UserAccount  string `xorm:"not null"`
		UserPassword string `xorm:"not null"`
	}{
		0, "user1", "passwd1",
	}
	engine.Table("user_account").CreateTable(&uAcc)
	engine.Table("user_account").Insert(&uAcc)

	uInfo := struct {
		UserId       uint32 `xorm:"pk unique not null"`
		UserName     string `xorm:"not null"`
		UserNickName string `xorm:"not null"`
		UserGender   string `xorm:"not null"`
	}{
		uAcc.UserId, "user1-name", "user1-nick-name", "man",
	}
	engine.Table("user_info").CreateTable(&uInfo)
	engine.Table("user_info").CreateUniques(&uInfo)
	engine.Table("user_info").Insert(&uInfo)

	eInst.AddClusterTable(&entity.IcClusterTable{
		ClassId: eClass.ClassId, ClusterName: "user account",
		ClusterDesc: "user account desc", ClusterTableName: "user_account",
		IsPrimary: true, TenantId: 1, Status: 1,
	})
	eInst.AddClusterTable(&entity.IcClusterTable{
		ClassId: eClass.ClassId, ClusterName: "user info",
		ClusterDesc: "user info desc", ClusterTableName: "user_info",
		IsPrimary: false, TenantId: 1, Status: 1,
	})
}
