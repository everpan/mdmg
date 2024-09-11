package entity

import (
	"github.com/everpan/mdmg/pkg/log"
	"github.com/everpan/mdmg/pkg/store"
	"go.uber.org/zap"
	"math"
	"xorm.io/xorm"
)

type IcEntityBase struct {
	EntityId   uint64 `json:"entity_id" xorm:"autoincr pk notnull"`
	ClassId    uint32 `xorm:"unique(cls_ukey) index"`
	EntityUKey string `json:"entity_ukey" xorm:"entity_ukey unique(cls_ukey)"`
}

// IcEntityClass 实体类划分属性簇；
// 尽可能每个簇属性放一个表内，通过 entity_id 关联
type IcEntityClass struct {
	ClassId      uint32           `xorm:"autoincr pk notnull"`
	ClassName    string           `xorm:"unique"`
	ClassDesc    string           `xorm:"text"` // 关于实体的描述信息
	ClusterInfos []*IcClusterInfo `xorm:"text"` // 属性表，第一个为主属性表
}

type IcClusterInfo struct {
	ClusterId        uint32 // 簇表名
	ClusterName      string // 簇名
	ClusterDesc      string `xorm:"text"`   // 簇名
	ClusterTableName string `xorm:"unique"` // unique 簇表名
}

var (
	engine           *xorm.Engine
	logger           *zap.Logger
	entityClassCache = store.OneLevelMap[uint32, *IcEntityClass]{}
)

func init() {
	logger = log.GetLogger()
}

func InitTable() {
	engine.CreateTables(&IcEntityBase{}, &IcEntityClass{}, &IcClusterInfo{})
	engine.CreateUniques(&IcEntityBase{})
	engine.CreateUniques(&IcEntityClass{})
	engine.CreateUniques(&IcClusterInfo{})
}

func SetEngine(eng *xorm.Engine) {
	engine = eng
}

func GetEntityClass(classId uint32) *IcEntityClass {
	e, ok := entityClassCache.Get(classId)
	if ok {
		return e
	}
	e = &IcEntityClass{ClassId: classId}
	ok, err := engine.Get(e)
	if err != nil {
		logger.Error("GetEntityClass", zap.Uint32("classId", classId), zap.Error(err))
		return nil
	}
	entityClassCache.Set(classId, e)
	return e
}

// AcquireEntityBase 以`ClassId,EntityUKey`进行插入或者查询对应的实体基信息
// 当 EntityId 为 0，即尝试插入
// 当 插入失败，则改为查询，如果查询不到，则返回 nil
func AcquireEntityBase(base *IcEntityBase) *IcEntityBase {
	if base.EntityId > 0 {
		// query
		engine.Get(base)
		return base
	} else {
		_, err := engine.Insert(base)
		if err != nil {
			logger.Warn("AcquireEntityBase insert entity.EntityId !=0,use QueryEntityBaseByClsIdAndUKey to query ",
				zap.Any("before entity:", base), zap.Error(err))
			ok, err := engine.Get(base)
			if err != nil {
				logger.Error("AcquireEntityBase query", zap.Any("before entity:", base), zap.Error(err))
			}
			if ok {
				logger.Warn("AcquireEntityBase query", zap.Any("queried entity:", base))
				return base
			}
		}
	}
	return nil
}

func QueryEntityBaseByClsIdAndUKey(clsId uint32, ukey string) *IcEntityBase {
	// // clsId=0的情况下 导致只会用EntityUKey去查询
	if clsId == 0 {
		clsId = math.MaxUint32
	}
	e := &IcEntityBase{
		ClassId:    clsId,
		EntityUKey: ukey,
	}
	ok, err := engine.Get(e)
	if ok {
		return e
	}
	if err != nil {
		logger.Warn("QueryEntityBaseByClsIdAndUKey", zap.Any("entity", e), zap.Error(err))
	}
	return nil
}
