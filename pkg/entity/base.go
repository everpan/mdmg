package entity

import (
	"errors"
	"fmt"
	"github.com/everpan/mdmg/pkg/log"
	"github.com/everpan/mdmg/pkg/store"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

// IcEntityClass 实体类； 划分属性簇； 注册管理，便于业务灵活定义
// 尽可能每个簇属性放一个表内，通过 entity_id 关联
type IcEntityClass struct {
	ClassId        uint32 `xorm:"autoincr pk notnull"`
	ClassName      string `xorm:"unique"`
	ClassDesc      string `xorm:"text"` // 关于实体的描述信息
	EntityPKColumn string `json:"entity_pk_column" xorm:"entity_pk_column index"`
	EntityUKColumn string `json:"entity_uk_column" xorm:"entity_uk_column index"` // 实体主键列名;统一实体的列类型为uint64，可以采用数据库自增
	// EntityPrimaryTable string           `xorm:"entity_primary_table unique"`
	ClusterColumns []*IcClusterColumn `xorm:"text"` // 属性表，第一个为主属性表; 所以的簇属性必需包含与`EntityPKColumn`同名的主键字段
}

type IcClusterColumn struct {
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
	engine.CreateTables(&IcEntityClass{})
	engine.CreateUniques(&IcEntityClass{})
	engine.CreateIndexes(&IcEntityClass{})

	engine.CreateTables(&IcClusterColumn{})
	engine.CreateUniques(&IcClusterColumn{})

}

func SetEngine(eng *xorm.Engine) {
	engine = eng
}

func insertNewEntityClass(ec *IcEntityClass) error {
	_, err := engine.Insert(ec)
	return err
}

func RegisterEntityClass(ec *IcEntityClass) *IcEntityClass {
	if ec.ClassId == 0 {
		err := insertNewEntityClass(ec)
		if err != nil {
			logger.Error("Failed to insert new entity class", zap.Error(err))
		}
		entityClassCache.Set(ec.ClassId, ec)
		return ec
	}

	e, ok := entityClassCache.Get(ec.ClassId)
	if ok {
		return e
	}
	return nil
}

func GetEntityClass(classId uint32) (*IcEntityClass, error) {
	if classId == 0 {
		return nil, errors.New("classId is 0")
	}
	e, ok := entityClassCache.Get(classId)
	if ok {
		return e, nil
	}
	ec := &IcEntityClass{ClassId: classId}
	var err error
	ok, err = engine.Get(ec)
	if err != nil {
		logger.Error("Failed to get entity class", zap.Error(err))
		return nil, err
	}
	if ok {
		entityClassCache.Set(ec.ClassId, ec)
		return ec, nil
	}
	return nil, errors.New(fmt.Sprintf("entity classId:%d not found", classId))
}
