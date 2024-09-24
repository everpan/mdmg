package entity

import (
	"errors"
	"fmt"
	"github.com/everpan/mdmg/dsl"
	"github.com/everpan/mdmg/pkg/base/log"
	"github.com/everpan/mdmg/pkg/base/store"
	"strconv"
	"strings"

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
	TenantId       uint32 `xorm:"index"`
	// EntityUKColumn string `json:"entity_uk_column" xorm:"entity_uk_column index"` // 实体主键列名;统一实体的列类型为uint64，可以采用数据库自增
	// EntityPrimaryTable string           `xorm:"entity_primary_table unique"`
	// ClusterIdList []uint32 `json:"cluster_id_list,omitempty" xorm:"cluster_id_list text default ''"` // 属性表，第一个为主属性表; 所以的簇属性必需包含与`EntityPKColumn`同名的主键字段
}

// IcClusterTable 簇表信息
// 多个簇表通过left join 进行查询，根据需求定制
// classId : clusterId = 1 : n
type IcClusterTable struct {
	ClassId          uint32 `xorm:"index"`               // 所属实体类
	ClusterId        uint32 `xorm:"pk autoincr notnull"` // 簇表名
	ClusterName      string // 簇名
	ClusterDesc      string `xorm:"text"`   // 簇描述
	ClusterTableName string `xorm:"unique"` // unique 簇表名； 至少包含EntityPKColumn
	IsPrimary        bool   `xorm:"bool"`   // 是否是主簇，主簇的key通常是自增
	TenantId         uint32 `xorm:"index"`
	Status           int32  `xorm:"index""` // 状态
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
//
//	ClusterId == 0  簇类为新
func (ctx *Context) AddClusterTable(ct *IcClusterTable) error {
	if ct.ClassId == 0 {
		return errors.New("classId is 0")
	}

	ec, err := ctx.GetEntityClass(ct.ClassId)
	if nil != err || ec == nil {
		return fmt.Errorf("entity classId:%d not found, err: %v ,entity:%v",
			ct.ClassId, err, ec)
	}
	return ctx.AddClusterTableWithoutCheckClassId(ct)
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

// CreateViewTable 以主表将所有簇表构建成view试图
// force 是否强制删除重建
func (ctx *Context) CreateViewTable(classId uint32, force bool) error {
	cTables, err := ctx.GetClusterTables(classId)
	if nil != err {
		return err
	}
	var primaryTable string
	var tNames []string

	for _, table := range cTables {
		if table.IsPrimary {
			primaryTable = table.ClusterTableName
			continue
		}
		tNames = append(tNames, table.ClusterTableName)
	}
	if len(primaryTable) == 0 {
		return fmt.Errorf("primary table is empty of classId:%d", classId)
	}

	viewTableName := fmt.Sprintf(`v_cls_%s`, primaryTable)
	_, err = ctx.engine.Exec("drop view if exists " + viewTableName)
	if err != nil {
		return err
	}
	if len(tNames) == 0 { // 没有非簇表
		_, err = ctx.engine.Exec("create view " + viewTableName)
	} else {
		_, err = ctx.engine.Exec("creat view " + viewTableName)
	}

	return nil
}

type void struct{}

// filterRepeatedColumns 获取簇表中相同列名
func filterRepeatedColumns(primaryTables *dsl.Meta, clusterTables []*dsl.Meta) map[string]void {
	var tmp = make(map[string]int)
	if primaryTables != nil {
		for _, col := range primaryTables.Columns {
			tmp[col.Name] += 1
		}
	}

	for _, clusterTable := range clusterTables {
		for _, col := range clusterTable.Columns {
			tmp[col.Name] += 1
		}
	}
	var result = make(map[string]void)
	var empty = void{}
	for colName, cnt := range tmp {
		if cnt > 1 {
			result[colName] = empty
		}
	}
	return result
}

// GenerateSelectColumnsSQL 生成sql选项，select t0.a t0.b ...
func GenerateSelectColumnsSQL(primaryTables *dsl.Meta, clusterTables []*dsl.Meta) (string, error) {
	var sb strings.Builder
	_, err := GenerateSelectColumnsSQLBuilder(&sb, primaryTables, clusterTables)
	return sb.String(), err
}

// GenerateSelectColumnsSQLBuilder 将
func GenerateSelectColumnsSQLBuilder(sb *strings.Builder, primaryTables *dsl.Meta, clusterTables []*dsl.Meta) (string, error) {
	var key string
	for _, col := range primaryTables.Columns {
		if col.IsPrimaryKey {
			key = col.Name
			break
		}
	}

	if key == "" {
		return key, fmt.Errorf("not fount primary key in table %s", primaryTables.Table.Name)
	}
	// allTables := append(clusterTables, primaryTables)
	repeatedCols := filterRepeatedColumns(primaryTables, clusterTables)
	var exist bool
	sb.WriteString("select t0.")
	sb.WriteString(key)
	for _, col := range primaryTables.Columns {
		if col.IsPrimaryKey {
			continue
		}
		sb.WriteString(", t0.")
		sb.WriteString(col.Name)
	}

	if clusterTables != nil {
		sb.WriteString(",\n")
	}

	for i, clusterTable := range clusterTables {
		var tIdx = strconv.Itoa(i + 1)
		var aliasTableDot = "t" + tIdx + "."
		for j, col := range clusterTable.Columns {
			if col.Name == key {
				continue
			}
			sb.WriteString(aliasTableDot)
			sb.WriteString(col.Name)
			_, exist = repeatedCols[col.Name]
			if exist {
				var aliasCol = col.Name + "_" + tIdx
				sb.WriteString(" as ")
				sb.WriteString(aliasCol)
			}
			if j < len(clusterTable.Columns)-1 {
				sb.WriteString(", ")
			}
		}
		if i < len(clusterTables)-1 {
			sb.WriteString(",\n")
		}
	}
	return key, nil
}

// GenerateLeftJoinSQL 将多个簇表left join起来
func GenerateLeftJoinSQL(clusterTables []*dsl.Meta, key string) string {
	var sb strings.Builder
	GenerateLeftJoinConditionSQLBuilder(&sb, clusterTables, key)
	return sb.String()
}

func GenerateLeftJoinConditionSQLBuilder(sb *strings.Builder, clusterTables []*dsl.Meta, key string) {
	if nil == clusterTables {
		return
	}
	for i, table := range clusterTables {
		alias := strconv.Itoa(i + 1)
		if i == 0 {
			sb.WriteString("left join ")
		} else {
			sb.WriteString("\nleft join ")
		}
		sb.WriteString(table.Table.Name)
		sb.WriteString(" as t")
		sb.WriteString(alias)
		sb.WriteString(" on t0")
		sb.WriteString(".")
		sb.WriteString(key)
		sb.WriteString(" = t")
		sb.WriteString(alias)
		sb.WriteString(".")
		sb.WriteString(key)
	}
}

func GenerateJoinTableSQL(primaryTables *dsl.Meta, clusterTables []*dsl.Meta) (string, error) {
	var sb strings.Builder
	key, err := GenerateSelectColumnsSQLBuilder(&sb, primaryTables, clusterTables)
	if nil != err {
		return "", err
	}
	sb.WriteString("\nfrom ")
	sb.WriteString(primaryTables.Table.Name)
	sb.WriteString(" as t0")
	if clusterTables != nil {
		sb.WriteString("\n")
	}
	GenerateLeftJoinConditionSQLBuilder(&sb, clusterTables, key)
	return sb.String(), nil
}

func FilterPrimaryClusterTable(tables []*IcClusterTable) *IcClusterTable {
	//if tables == nil {
	//	return nil
	//}
	for _, table := range tables {
		if table.IsPrimary {
			return table
		}
	}
	return nil
}
