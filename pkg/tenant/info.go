package tenant

import (
	"errors"
	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/everpan/mdmg/pkg/log"
	"github.com/everpan/mdmg/pkg/store"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

const DefaultGuidNamespace = "11111111-1111-1111-1111-111111111111"

type Info struct {
	Idx           uint32 `json:"idx" xorm:"tenant_id pk autoincr"`
	SId           string `json:"sid" xorm:"char(36) notnull unique"`
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
	cache               = store.OneLevelMap[string, *Info]{}
	namespace           guid.GUID
	defaultSystemEngine *xorm.Engine
	engineCache         = store.OneLevelMap[string, *xorm.Engine]{}
)

var (
	DefaultInfo = Info{
		Idx:       1,
		SId:       DefaultGuidNamespace,
		En:        "default_test",
		Cn:        "默认租户",
		IsTestEnv: true,
	}
	DefaultHostInfo = Info{
		Idx:    2,
		SId:    "22222222-2222-2222-2222-222222222222",
		En:     "host",
		Cn:     "运营商",
		IsHost: true,
	}
)

const InfoTableName = "ic_tenant_info"

func SetEngine(e *xorm.Engine) {
	defaultSystemEngine = e
}

func init() {
	namespace, _ = guid.FromString(DefaultGuidNamespace)
}

func NewTenantInfo(sid string, en string, cn string, extension string, isTest bool) *Info {
	id, _ := guid.NewV5(namespace, []byte(sid))
	return &Info{
		Idx: 0,
		SId: id.String(),
		En:  en, Cn: cn, IsTestEnv: isTest,
		Extension: extension,
	}
}

func (tenant *Info) SaveTenantInfo() error {
	var err error
	if tenant.Idx == 0 {
		_, err = defaultSystemEngine.Insert(tenant)
	} else {
		_, err = defaultSystemEngine.Update(tenant)
	}
	return err
}

func AcquireTenantInfoBySid(sid string) (*Info, error) {
	info, ok := cache.Get(sid)
	if ok {
		return info, nil
	}
	info = &Info{SId: sid}
	var err error
	if nil == defaultSystemEngine {
		return nil, errors.New("default system engine is null")
	}
	ok, err = defaultSystemEngine.Table(InfoTableName).Get(info)
	if ok {
		cache.Set(sid, info)
		return info, nil
	}
	if nil != err {
		log.GetLogger().Error("query tenant", zap.String("sid", sid), zap.Error(err))
	}
	return nil, err
}

func AcquireEngine(info *Info) (*xorm.Engine, error) {
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
