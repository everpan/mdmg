package entity

import (
	"os"
	"testing"

	"xorm.io/xorm"
)

type UserEmailAccout struct {
	UserId    uint32 `xorm:"autoincr pk notnull"`
	UserEmail string `xorm:"unique"`
	Password  string
}
type UserNickNameAccoutn struct {
	UserId   uint32 `xorm:"autoincr pk notnull"`
	NickName string `xorm:"unique"`
	Password string
}

func TestRegister(t *testing.T) {
	dbFile := "./register_test.db"
	_ = os.Remove(dbFile)
	engine, err := xorm.NewEngine("sqlite3", dbFile)
	if nil != err {
		t.Fatal(err)
	}
	ctx := NewContext(engine)
	ctx.InitTable()

	engine.CreateTables(&UserEmailAccout{})
	engine.CreateUniques(&UserEmailAccout{})
	engine.CreateTables(&UserNickNameAccoutn{})
	engine.CreateUniques(&UserNickNameAccoutn{})

	emailTName := engine.TableName(&UserEmailAccout{})
	nickTName := engine.TableName(&UserNickNameAccoutn{})
	t.Logf("NEW TABLE: %s %s", emailTName, nickTName)
}
