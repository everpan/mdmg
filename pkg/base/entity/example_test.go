package entity

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"xorm.io/xorm"
)

type UserEmailAccount struct {
	UserId       uint32 `xorm:"autoincr pk notnull"`
	UserEmail    string `xorm:"unique"`
	UserPassword string
}
type UserNickNameAccount struct {
	UserId       uint32 `xorm:"index notnull"`
	NickName     string `xorm:"unique"`
	UserPassword string
}

func TestRegister(t *testing.T) {
	dbFile := "./register_test.db"
	_ = os.Remove(dbFile)
	engine, err := xorm.NewEngine("sqlite3", dbFile)
	engine.ShowSQL(true)
	if nil != err {
		t.Fatal(err)
	}
	ctx := NewContext(engine)
	ctx.InitTable()
	uea := &UserEmailAccount{}
	una := &UserNickNameAccount{}
	engine.CreateTables(uea)
	engine.CreateUniques(uea)
	engine.CreateTables(una)
	engine.CreateIndexes(una)
	engine.CreateUniques(una)

	emailTName := engine.TableName(uea)
	nickTName := engine.TableName(una)
	t.Logf("NEW TABLE: %s %s", emailTName, nickTName)
	//
	ec := &IcEntityClass{
		ClassName: "UserAccount",
		ClassDesc: "用户账号",
		PkColumn:  "user_id",
	}
	// 1. 注册实体类 UserAccount // 注意大小写
	ctx.RegisterEntityClass(ec)
	fmt.Sprintf("register entity cliass %v", ec)
	assert.Greater(t, ec.ClassId, uint32(0))

	// class id must gt 0
	ct0 := &IcClusterTable{}
	err = ctx.AddClusterTable(ct0)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "classId is 0")

	ct1 := &IcClusterTable{
		ClassId:          ec.ClassId,
		ClusterName:      "用户邮件账户",
		ClusterDesc:      "用户邮件账户，用于登录可选",
		ClusterTableName: emailTName,
		IsPrimary:        true,
	}
	_ = ctx.AddClusterTable(ct1)
	assert.Greater(t, ct1.ClusterId, uint32(0))
	ct2 := &IcClusterTable{
		ClassId:          ec.ClassId,
		ClusterName:      "用户昵称账户",
		ClusterDesc:      "用户昵称账户，用于登录可选",
		ClusterTableName: nickTName,
	}
	err = ctx.AddClusterTable(ct2)
	assert.Nil(t, err)
	assert.Greater(t, ct2.ClusterId, ct1.ClusterId)

	tables, err := ctx.GetClusterTables(ec.ClassId)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tables))
	primaryTable := FilterPrimaryClusterTable(nil)
	assert.Nil(t, primaryTable)
	primaryTable = FilterPrimaryClusterTable(tables)
	assert.NotNil(t, primaryTable)
}
