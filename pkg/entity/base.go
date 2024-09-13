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
	// EntityUKColumn string `json:"entity_uk_column" xorm:"entity_uk_column index"` // 实体主键列名;统一实体的列类型为uint64，可以采用数据库自增
	// EntityPrimaryTable string           `xorm:"entity_primary_table unique"`
	// ClusterIdList []uint32 `json:"cluster_id_list,omitempty" xorm:"cluster_id_list text default ''"` // 属性表，第一个为主属性表; 所以的簇属性必需包含与`EntityPKColumn`同名的主键字段
}

type IcClusterTable struct {
	ClusterId        uint32 `xorm:"pk autoincr notnull"` // 簇表名
	ClusterName      string // 簇名
	ClassId          uint32 `xorm:"index"`  // 所属实体类
	ClusterDesc      string `xorm:"text"`   // 簇描述
	ClusterTableName string `xorm:"unique"` // unique 簇表名； 至少包含EntityPKColumn
	IsPrimary        bool   `xorm:"bool"`   // 是否是主簇，主簇的key通常是自增
}

var (
	logger *zap.Logger
)

func init() {
	logger = log.GetLogger()
}

type Context struct {
	engine           *xorm.Engine
	entityClassCache store.OneLevelMap[uint32, *IcEntityClass]
}

func NewContext(engine *xorm.Engine) *Context {
	return &Context{
		engine:           engine,
		entityClassCache: store.OneLevelMap[uint32, *IcEntityClass]{},
	}
}

func (ctx *Context) InitTable() {
	ec := &IcEntityClass{}
	ct := &IcClusterTable{}
	engine := ctx.engine

	_ = engine.CreateTables(ec, ct)
	_ = engine.CreateUniques(ec)
	_ = engine.CreateIndexes(ec)

	_ = engine.CreateUniques(ct)
	_ = engine.CreateIndexes(ct)
}

func (ctx *Context) SetEngine(eng *xorm.Engine) {
	ctx.engine = eng
}

func (ctx *Context) insertNewEntityClass(ec *IcEntityClass) error {
	_, err := ctx.engine.Insert(ec)
	return err
}

func (ctx *Context) RegisterEntityClass(ec *IcEntityClass) (*IcEntityClass, error) {
	if ec.ClassId == 0 {
		err := ctx.insertNewEntityClass(ec)
		if err != nil {
			logger.Error("Failed to insert new entity class", zap.Error(err))
			return nil, err
		}
		ctx.entityClassCache.Set(ec.ClassId, ec)
		return ec, nil
	}
	return ctx.GetEntityClass(ec.ClassId)
}

// RegisterClassName 注册实体类名，其他信息后续补充，否则不能工作；主要简化工作

func (ctx *Context) RegisterClassName(className string) (*IcEntityClass, error) {
	ec := &IcEntityClass{ClassName: className}
	return ctx.RegisterEntityClass(ec)
}

func (ctx *Context) GetEntityClass(classId uint32) (*IcEntityClass, error) {
	if classId == 0 {
		return nil, errors.New("classId is 0")
	}
	e, ok := ctx.entityClassCache.Get(classId)
	if ok {
		return e, nil
	}
	ec := &IcEntityClass{ClassId: classId}
	var err error
	ok, err = ctx.engine.Get(ec)
	if err != nil {
		logger.Error("Failed to get entity class", zap.Error(err))
		return nil, err
	}
	if ok {
		ctx.entityClassCache.Set(ec.ClassId, ec)
		return ec, nil
	}
	return nil, fmt.Errorf("entity classId:%d not found", classId)
}

func (ctx *Context) GetEntityClassByName(className string) (*IcEntityClass, error) {
	ec := &IcEntityClass{}
	ok, err := ctx.engine.Where("class_name = ?", className).Get(ec)
	if err != nil {
		logger.Error("Failed to get entity class", zap.Error(err))
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("entity className:%s not found", className)
	}
	return ec, nil
}

// AddClusterTableWithoutCheckClassId classId == 0 的情况下，注册簇属性表
func (ctx *Context) AddClusterTableWithoutCheckClassId(ct *IcClusterTable) error {
	if ct.ClusterId != 0 {
		return fmt.Errorf("clusterId:%d is not 0, pls use GetClusterTable to get details", ct.ClusterId)
	}
	var (
		err error
	)
	_, err = ctx.engine.Insert(ct)
	return err
}

// AddClusterTable 增加簇表
// 条件 ： classId > 0     存在实体类
//		  ClusterId == 0  簇类为新

func (ctx *Context) AddClusterTable(ct *IcClusterTable) error {
	if ct.ClassId == 0 {
		return errors.New("classId is 0")
	}

	ec, err1 := ctx.GetEntityClass(ct.ClassId)
	if nil != err1 {
		return err1
	}
	if ec == nil {
		return fmt.Errorf("entity classId:%d not found", ct.ClassId)
	}
	return nil
}

func (ctx *Context) GetClusterTables(classId uint32) ([]*IcClusterTable, error) {
	tables := make([]*IcClusterTable, 0)
	err := ctx.engine.Where("class_id = ?", classId).Find(&tables)
	return tables, err
}

func (ctx *Context) GetPrimaryClusterTable(classId uint32) (*IcClusterTable, error) {
	table := &IcClusterTable{ClassId: classId, IsPrimary: true}
	ok, err := ctx.engine.Get(table)
	if err != nil {
		return nil, err
	}
	if ok {
		return table, nil
	}
	return nil, nil
}

func FilterPrimaryClusterTable(tables []*IcClusterTable) *IcClusterTable {
	if tables == nil {
		return nil
	}
	for _, table := range tables {
		if table.IsPrimary {
			return table
		}
	}
	return nil
}
