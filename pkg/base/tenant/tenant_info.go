package tenant

import (
	"errors"
	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/everpan/mdmg/pkg/base/log"
	"github.com/everpan/mdmg/pkg/base/store"
	"github.com/everpan/mdmg/pkg/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

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
	cache     = store.OneLevelMap[string, *IcTenantInfo]{}
	namespace guid.GUID
	sysEngine *xorm.Engine // 租户管理为最高权限，运营商才可
)

var (
	logger         = log.GetLogger()
	Enable         bool
	TestTenantInfo *IcTenantInfo
	HostTenantInfo *IcTenantInfo
)

func init() {
	var (
		hostSid = "22222222-2222-2222-2222-222222222222"
		TestSid = "11111111-1111-1111-1111-111111111111"
	)
	namespace, _ = guid.FromString(hostSid)
	TestTenantInfo = NewTenantInfo(0, TestSid, "test-tenant", "测试", "", true, true)
	HostTenantInfo = NewTenantInfo(0, hostSid, "host-tenant", "运营商", "", false, true)
	TestTenantInfo.Driver = "sqlite3"
	TestTenantInfo.ConnectString = "ic-default.db"
	HostTenantInfo.Driver = "sqlite3"
	HostTenantInfo.ConnectString = "ic-default.db"

	viper.SetDefault("tenant.db-driver", "sqlite3")
	viper.SetDefault("tenant.db-connect", "./ic-tenant.db")
	viper.SetDefault("tenant.enable", false)
	config.RegisterReloadViperFunc(updateConfig)
}

func initTable(engine *xorm.Engine) {
	t := &IcTenantInfo{}
	err := engine.CreateTables(t)
	if nil != err {
		_ = engine.CreateUniques(t)
	}
}

func updateConfig() error {
	Enable = viper.GetBool("tenant.enable")
	if Enable {
		driver := viper.GetString("tenant.db-driver")
		connectString := viper.GetString("tenant.db-connect")

		logger.Info("tenant admin engine", zap.String("driver", driver))

		engine, err := config.AcquireEngine(driver, connectString)
		if err != nil {
			panic(err)
		}
		sysEngine = engine
		initTable(engine)
		TestTenantInfo = syncInfoFromDB(TestTenantInfo)
		HostTenantInfo = syncInfoFromDB(HostTenantInfo)
	}
	logger.Info("tenant-config", zap.Any("host-sid", HostTenantInfo),
		zap.Any("test-sid", TestTenantInfo), zap.Bool("enable", Enable))
	return nil
}

func syncInfoFromDB(info *IcTenantInfo) *IcTenantInfo {
	ok, _ := sysEngine.Exist(info)
	if !ok {
		info.Save()
	} else {
		info, _ = AcquireInfoBySid(TestTenantInfo.SId)
	}
	return info
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

func (t *IcTenantInfo) IsTest() bool {
	return t.IsTestEnv
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
