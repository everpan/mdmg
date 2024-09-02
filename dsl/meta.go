package dsl

import (
	"context"
	"encoding/json"
	"strings"
	"xorm.io/xorm"
	"xorm.io/xorm/core"
	"xorm.io/xorm/dialects"
	"xorm.io/xorm/schemas"
)

type Table struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Charset string `json:"charset,omitempty"`
}

type Column struct {
	Name            string   `json:"name"`
	Label           string   `json:"label,omitempty"` // 当从数据库中reverse的时候，丢失了label信息
	Type            string   `json:"type"`
	Length1         int64    `json:"length1,omitempty"`
	Length2         int64    `json:"length2,omitempty"`
	Default         string   `json:"default,omitempty"`
	Options         []string `json:"options,omitempty"`
	IsPrimaryKey    bool     `json:"primary_key,omitempty"`
	IsAutoIncrement bool     `json:"auto_increment,omitempty"`
	Index           bool     `json:"index,omitempty"`
	Unique          bool     `json:"unique,omitempty"`
	Nullable        bool     `json:"nullable,omitempty"`
	Comment         string   `json:"comment,omitempty"`
}

func (c *Column) GetSQLType() *schemas.SQLType {
	return &schemas.SQLType{
		Name:           strings.ToUpper(c.Type),
		DefaultLength:  c.Length1,
		DefaultLength2: c.Length2,
	}
}

func (c *Column) ConvertXormColumn(db *core.DB) *schemas.Column {
	sqlType := *c.GetSQLType()
	name := db.Mapper.Obj2Table(c.Name)
	col := schemas.NewColumn(name, name, sqlType, c.Length1, c.Length2, c.Nullable)
	col.Comment = c.Comment
	col.IsPrimaryKey = c.IsPrimaryKey

	if c.Options != nil && len(c.Options) > 0 {
		if sqlType.Name == schemas.Enum {
			for i, option := range c.Options {
				col.EnumOptions[option] = i
			}
		}
	}

	if len(c.Default) > 0 {
		col.DefaultIsEmpty = false
		col.Default = c.Default
	}

	c.parseIDAssignXormColumn(col)

	return col
}

// ID,ID32,ID64 为特殊tag，todo 反向转换
func (c *Column) parseIDAssignXormColumn(col *schemas.Column) {
	// ID -> bigint unsigned not null primary key auto_increment
	// ID32 -> int unsigned not null primary key auto_increment
	if col.SQLType.Name == "ID" || col.SQLType.Name == "ID64" {
		// sqlType.Name = schemas.BigInt // sqlType is temp
		col.SQLType.Name = schemas.UnsignedBigInt
		col.IsPrimaryKey = true
		col.Nullable = false
		col.IsAutoIncrement = true
	} else if col.SQLType.Name == "ID32" {
		col.SQLType.Name = schemas.UnsignedInt
		col.IsPrimaryKey = true
		col.Nullable = false
		col.IsAutoIncrement = true
	}
}

type Meta struct {
	Table   Table     `json:"table"`
	Columns []*Column `json:"columns"`
}

func (meta *Meta) CreateTableSchema(db *core.DB, dialect dialects.Dialect) *schemas.Table {
	table := schemas.NewEmptyTable()
	table.Name = db.Mapper.Obj2Table(meta.Table.Name)
	table.Comment = meta.Table.Comment
	// columns
	for _, c := range meta.Columns {
		col := c.ConvertXormColumn(db)
		table.AddColumn(col)
	}
	if len(strings.TrimSpace(meta.Table.Charset)) == 0 {
		table.Charset = "utf8mb4"
	}
	return table
}

func (meta *Meta) CreateTableSQL(db *core.DB, dialect dialects.Dialect) (sqlStr string, err error) {
	table := meta.CreateTableSchema(db, dialect)
	sqlStr, _, err = dialect.CreateTableSQL(context.Background(), db, table, "")
	return
}

func (meta *Meta) CreateTable(eng *xorm.Engine) error {
	sqlStr, err := meta.CreateTableSQL(eng.DB(), eng.Dialect())
	if err != nil {
		return err
	}
	_, err = eng.Exec(sqlStr)
	return err
}

func (meta *Meta) FromSchemaTable(t *schemas.Table) {
	meta.Table.Name = t.Name
	meta.Table.Comment = t.Comment
	meta.Table.Charset = t.Charset
	if meta.Columns == nil {
		meta.Columns = make([]*Column, 0, len(t.Columns()))
	}
	for _, c := range t.Columns() {
		col := (*XormColumn)(c).ConvertDslColumn()
		meta.Columns = append(meta.Columns, col)
	}
}

// Json  化整个meta
func (meta *Meta) Json(indent bool) (string, error) {
	var (
		d   []byte
		err error
	)
	if indent {
		d, err = json.MarshalIndent(meta, "", " ")

	} else {
		d, err = json.Marshal(meta)
	}
	return string(d), err
}

type XormColumn schemas.Column

func (xc *XormColumn) ConvertDslColumn() *Column {
	c := &Column{
		Name:            xc.Name,
		Type:            strings.ToLower(xc.SQLType.Name),
		Length1:         xc.SQLType.DefaultLength,
		Length2:         xc.SQLType.DefaultLength2,
		Default:         xc.Default,
		IsPrimaryKey:    xc.IsPrimaryKey,
		IsAutoIncrement: xc.IsAutoIncrement,
		Nullable:        xc.Nullable,
		Comment:         xc.Comment,
	}

	if len(xc.EnumOptions) > 0 {
		c.Options = make([]string, 0, len(xc.EnumOptions))
		for option := range xc.EnumOptions {
			c.Options = append(c.Options, option)
		}
	}
	return c
}
