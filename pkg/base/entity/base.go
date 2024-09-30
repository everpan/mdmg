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
	ClassId   uint32 `json:"class_id" xorm:"autoincr pk notnull"`
	ClassName string `json:"name" xorm:"unique"`
	ClassDesc string `json:"class_desc" xorm:"text"` // 关于实体的描述信息
	PkColumn  string `json:"pk_column" xorm:"pk_column index"`
	TenantId  uint32 `json:"tenant_id" xorm:"index"`
	Status    int32  `json:"status" xorm:"index""` // 状态
	// EntityUKColumn string `json:"entity_uk_column" xorm:"entity_uk_column index"` // 实体主键列名;统一实体的列类型为uint64，可以采用数据库自增
	// EntityPrimaryTable string           `xorm:"entity_primary_table unique"`
	// ClusterIdList []uint32 `json:"cluster_id_list,omitempty" xorm:"cluster_id_list text default ''"` // 属性表，第一个为主属性表; 所以的簇属性必需包含与`PkColumn`同名的主键字段
}

// IcClusterTable 簇表信息
// 多个簇表通过left join 进行查询，根据需求定制
// classId : clusterId = 1 : n
type IcClusterTable struct {
	ClassId          uint32 `json:"class_id" xorm:"index"`                 // 所属实体类
	ClusterId        uint32 `json:"cluster_id" xorm:"pk autoincr notnull"` // 簇表名
	ClusterName      string `json:"name"`                                  // 簇名
	ClusterDesc      string `json:"cluster_desc" xorm:"text"`              // 簇描述
	ClusterTableName string `json:"table" xorm:"unique"`                   // unique 簇表名； 至少包含EntityPKColumn
	IsPrimary        bool   `json:"is_primary" xorm:"bool"`                // 是否是主簇，主簇的key通常是自增
	TenantId         uint32 `json:"tenant_id" xorm:"index"`
	Status           int32  `json:"status" xorm:"index"` // 状态
}
type IcEntityMeta struct {
	EntityClass   *IcEntityClass    `json:"entity_class"`
	ClusterTables []*IcClusterTable `json:"cluster_tables"`
}

var (
	logger *zap.Logger
)

func init() {
	logger = log.GetLogger()
}

type Context struct {
	engine           *xorm.Engine
	tenantId         uint32
	entityClassCache store.OneLevelMap[uint32, *IcEntityClass]
}

func NewContext(engine *xorm.Engine, tenantId uint32) *Context {
	return &Context{
		engine:           engine,
		tenantId:         tenantId,
		entityClassCache: store.OneLevelMap[uint32, *IcEntityClass]{},
	}
}

func (ctx *Context) InitTable(engine *xorm.Engine) error {
	ec := &IcEntityClass{}
	ct := &IcClusterTable{}

	err := engine.CreateTables(ec, ct)
	if err != nil {
		return err
	}
	_ = engine.CreateUniques(ec)
	_ = engine.CreateIndexes(ec)

	_ = engine.CreateUniques(ct)
	_ = engine.CreateIndexes(ct)
	return nil
}

func (ctx *Context) SetEngine(eng *xorm.Engine) {
	ctx.engine = eng
}
func (ctx *Context) SetTenantId(tenantId uint32) {
	ctx.tenantId = tenantId
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
	return ctx.GetEntityClassById(ec.ClassId)
}

// RegisterClassName 注册实体类名，其他信息后续补充，否则不能工作；主要简化工作
func (ctx *Context) RegisterClassName(className string) (*IcEntityClass, error) {
	ec := &IcEntityClass{ClassName: className}
	return ctx.RegisterEntityClass(ec)
}

func (ctx *Context) GetEntityClassById(classId uint32) (*IcEntityClass, error) {
	if classId == 0 {
		return nil, errors.New("classId is 0")
	}
	e, ok := ctx.entityClassCache.Get(classId)
	if ok {
		return e, nil
	}
	ec := &IcEntityClass{ClassId: classId, TenantId: ctx.tenantId}
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
	return nil, fmt.Errorf("entity classId: %d tenantId: %d not found", classId, ctx.tenantId)
}

func (ctx *Context) GetEntityClassByName(className string) (*IcEntityClass, error) {
	ec := &IcEntityClass{}
	ok, err := ctx.engine.Where("class_name = ? and tenant_id = ?", className, ctx.tenantId).Get(ec)
	if err != nil {
		logger.Error("Failed to get entity class", zap.Error(err))
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("entity className:%s tenantId:%d not found", className, ctx.tenantId)
	}
	return ec, nil
}

func (ctx *Context) DelEntityClassById(classId uint32) error {
	err := ctx.DelAllClusterTableByClassId(classId)
	if err != nil {
		return err
	}
	sql := "delete from ic_entity_class where class_id = ? and tenant_id = ?"
	_, err = ctx.engine.Exec(sql, classId, ctx.tenantId)
	ctx.entityClassCache.Release(classId)
	return err
}

func (ctx *Context) DelEntityClassByName(className string) error {
	entityClass, err := ctx.GetEntityClassByName(className)
	if nil != err {
		return err
	}
	return ctx.DelEntityClassById(entityClass.ClassId)
}

func (ctx *Context) DelClusterTableById(clusterId uint32) error {
	sql := "delete from ic_cluster_table where cluster_id = ? and tenant_id = ?"
	_, err := ctx.engine.Exec(sql, clusterId, ctx.tenantId)
	return err
}

func (ctx *Context) DelAllClusterTableByClassId(classId uint32) error {
	sql := "delete from ic_cluster_table where class_id = ? and tenant_id = ?"
	_, err := ctx.engine.Exec(sql, classId, ctx.tenantId)
	return err
}

func (ctx *Context) DelClusterTableByTableName(clusterTableName string) error {
	clusterTable, err := ctx.GetClusterTableByClusterName(clusterTableName)
	if nil != err {
		return err
	}
	return ctx.DelClusterTableById(clusterTable.ClusterId)
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

	ec, err := ctx.GetEntityClassById(ct.ClassId)
	if nil != err || ec == nil {
		return fmt.Errorf("entity classId:%d not found, err: %v ,entity:%v",
			ct.ClassId, err, ec)
	}
	return ctx.AddClusterTableWithoutCheckClassId(ct)
}

func (ctx *Context) GetClusterTablesByClassId(classId uint32) ([]*IcClusterTable, error) {
	tables := make([]*IcClusterTable, 0)
	err := ctx.engine.Where("class_id = ? and tenant_id = ?", classId, ctx.tenantId).Find(&tables)
	return tables, err
}

func (ctx *Context) GetPrimaryClusterTable(classId uint32) (*IcClusterTable, error) {
	table := &IcClusterTable{ClassId: classId, IsPrimary: true, TenantId: ctx.tenantId}
	ok, err := ctx.engine.Get(table)
	if err != nil {
		return nil, err
	}
	if ok {
		return table, nil
	}
	return nil, nil
}

func (ctx *Context) GetClusterTableByClusterName(tableName string) (*IcClusterTable, error) {
	table := &IcClusterTable{ClusterTableName: tableName, TenantId: ctx.tenantId}
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
	cTables, err := ctx.GetClusterTablesByClassId(classId)
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
