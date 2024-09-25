package tenant

import (
	"errors"
	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/everpan/mdmg/pkg/base/log"
	"github.com/everpan/mdmg/pkg/base/store"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

const DefaultGuidNamespace = "11111111-1111-1111-1111-111111111111"

type IcTenantInfo struct {
	Idx             uint32 `json:"idx" xorm:"tenant_id pk autoincr"`
	SId             string `json:"sid" xorm:"tenant_sid char(36) notnull unique"`
	EnName          string `json:"en" xorm:"en_name varchar(64) notnull unique"` // 英文名称，用于登录，区分租户
	CnName          string `json:"cn" xorm:"cn_name varchar(128) notnull"`       // 中文名称
	Driver          string `json:"driver" xorm:"varchar(128) notnull"`
	ConnectString   string `json:"connect_string" xorm:"varchar(1024) notnull"` // 数据库连接配置
	LastUpTime      int64  `json:"last_up_time" xorm:"-"`                       // 内存使用，为效率，未同步数据库；用于清理cache
	ExtensionConfig string `json:"extension_config" xorm:"text "`               // 通过json进行扩展属性
	IsTestEnv       bool   `json:"is_test_env" xorm:"boolean"`                  // 测试环境
	IsHost          bool   `json:"is_host" xorm:"-"`                            // 运营商账号，每个系统只有一个
}

var (
	cache       = store.OneLevelMap[string, *IcTenantInfo]{}
	namespace   guid.GUID
	sysEngine   *xorm.Engine // 租户管理为最高权限，运营商才可
	engineCache = store.OneLevelMap[string, *xorm.Engine]{}
)

var (
	DefaultInfo = NewTenantInfo(uint32(1),
		DefaultGuidNamespace, "default_test", "默认租户", "", false, true)
	DefaultHostInfo = NewTenantInfo(uint32(2),
		"22222222-2222-2222-2222-222222222222", "host", "运营商", "", true, true)
)

// SetSysEngine 专门用于系统管理，区别于应用db
func SetSysEngine(e *xorm.Engine) {
	sysEngine = e
	engineCache.Set(e.DataSourceName(), e)
}

func init() {
	namespace, _ = guid.FromString(DefaultGuidNamespace)
	DefaultInfo.SId = DefaultGuidNamespace
	DefaultHostInfo.SId = "22222222-2222-2222-2222-222222222222"
}

func NewTenantInfo(idx uint32, sid string, en string, cn string,
	extension string, isTest bool, isHost bool) *IcTenantInfo {
	gid, _ := guid.NewV5(namespace, []byte(sid))
	return &IcTenantInfo{
		Idx:    idx,
		SId:    gid.String(),
		EnName: en, CnName: cn, IsTestEnv: isTest, IsHost: isHost,
		ExtensionConfig: extension,
	}
}

func (t *IcTenantInfo) Save() error {
	var err error
	if t.Idx == 0 {
		_, err = sysEngine.Insert(t)
	} else {
		_, err = sysEngine.Update(t)
	}
	return err
}

func (t *IcTenantInfo) InitTable(engine *xorm.Engine) error {
	err := engine.CreateTables(t)
	if nil != err {
		return err
	}
	_ = engine.CreateUniques(t)
	return nil
}

func AcquireInfoBySid(sid string) (*IcTenantInfo, error) {
	// sid : string id
	info, ok := cache.Get(sid)
	if ok {
		return info, nil
	}
	info = &IcTenantInfo{SId: sid}
	var err error
	if nil == sysEngine {
		return nil, errors.New("default system sysEngine is null")
	}
	ok, err = sysEngine.Get(info)
	if ok {
		cache.Set(sid, info)
		return info, nil
	}
	if nil != err {
		log.GetLogger().Error("query tenant", zap.String("sid", sid), zap.Error(err))
	}
	return nil, err
}

// AcquireEngineForTenant 获取租户的db连接；
// 以connect string作为key
func AcquireEngineForTenant(t *IcTenantInfo) (*xorm.Engine, error) {
	driver, connStr := t.Driver, t.ConnectString
	eng, ok := engineCache.Get(connStr)
	if ok {
		return eng, nil
	}
	var err error
	eng, err = xorm.NewEngine(driver, connStr)
	if err != nil {
		log.GetLogger().Error("create sysEngine failed",
			zap.String("driver", driver), zap.String("connStr", connStr), zap.Error(err))
		return nil, err
	}
	engineCache.Set(connStr, eng)
	return eng, nil
}
