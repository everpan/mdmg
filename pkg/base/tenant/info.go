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
	Idx           uint32 `json:"idx" xorm:"tenant_id pk autoincr"`
	SId           string `json:"sid" xorm:"sid char(36) notnull unique"`
	En            string `json:"en" xorm:"varchar(64) notnull unique"` // 英文名称，用于登录，区分租户
	Cn            string `json:"cn" xorm:"varchar(128) notnull"`       // 中文名称
	Driver        string `json:"driver" xorm:"varchar(128) notnull"`
	ConnectString string `json:"connect_string" xorm:"varchar(1024) notnull"` // 数据库连接配置
	LastUpTime    int64  `json:"last_up_time" xorm:"-"`                       // 内存使用，为效率，未同步数据库；用于清理cache
	Extension     string `json:"extension" xorm:"text "`                      // 通过json进行扩展属性
	IsTestEnv     bool   `json:"is_test_env" xorm:"boolean"`                  //测试环境
	IsHost        bool   `json:"is_host" xorm:"-"`                            // 运营商账号，每个系统只有一个
}

var (
	cache               = store.OneLevelMap[string, *IcTenantInfo]{}
	namespace           guid.GUID
	defaultSystemEngine *xorm.Engine // 租户管理为最高权限，运营商才可
	engineCache         = store.OneLevelMap[string, *xorm.Engine]{}
)

var (
	DefaultInfo = NewTenantInfo(uint32(1),
		DefaultGuidNamespace, "default_test", "默认租户", "", false)
	DefaultHostInfo = NewTenantInfo(uint32(2),
		"22222222-2222-2222-2222-222222222222", "host", "运营商", "", true)
)

// SetSysEngine 专门用于系统管理，区别于应用db
func SetSysEngine(e *xorm.Engine) {
	defaultSystemEngine = e
}

func init() {
	namespace, _ = guid.FromString(DefaultGuidNamespace)
	DefaultInfo.SId = DefaultGuidNamespace
	DefaultHostInfo.SId = "22222222-2222-2222-2222-222222222222"
}

func NewTenantInfo(idx uint32, sid string, en string, cn string, extension string, isTest bool) *IcTenantInfo {
	gid, _ := guid.NewV5(namespace, []byte(sid))
	return &IcTenantInfo{
		Idx: idx,
		SId: gid.String(),
		En:  en, Cn: cn, IsTestEnv: isTest,
		Extension: extension,
	}
}

func (tenant *IcTenantInfo) SaveTenantInfo() error {
	var err error
	if tenant.Idx == 0 {
		_, err = defaultSystemEngine.Insert(tenant)
	} else {
		_, err = defaultSystemEngine.Update(tenant)
	}
	return err
}

func AcquireTenantInfoBySid(sid string) (*IcTenantInfo, error) {
	// sid : string id
	info, ok := cache.Get(sid)
	if ok {
		return info, nil
	}
	info = &IcTenantInfo{SId: sid}
	var err error
	if nil == defaultSystemEngine {
		return nil, errors.New("default system engine is null")
	}
	ok, err = defaultSystemEngine.Get(info)
	if ok {
		cache.Set(sid, info)
		return info, nil
	}
	if nil != err {
		log.GetLogger().Error("query tenant", zap.String("sid", sid), zap.Error(err))
	}
	return nil, err
}

func AcquireTenantEngine(info *IcTenantInfo) (*xorm.Engine, error) {
	driver, connStr := info.Driver, info.ConnectString
	eng, ok := engineCache.Get(connStr)
	if ok {
		return eng, nil
	}
	var err error
	eng, err = xorm.NewEngine(driver, connStr)
	if err != nil {
		log.GetLogger().Error("create defaultSystemEngine failed",
			zap.String("driver", driver), zap.String("connStr", connStr), zap.Error(err))
		return nil, err
	}
	engineCache.Set(connStr, eng)
	return eng, nil
}
