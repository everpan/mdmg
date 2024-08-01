package models

import (
	"reflect"
	"xorm.io/xorm/schemas"
)

type SQLType struct {
	Name           string `json:"sql_type_name"`
	DefaultLength  int64  `json:"default_length"`
	DefaultLength2 int64  `json:"default_length_2"`
}

type Column struct {
	ID              uint64   `xorm:"id pk autoincr" json:"id"`
	ModelID         uint64   `xorm:"model_id index unique('uk_model_name')" json:"model_id"`
	Name            string   `xorm:"col_name index unique('uk_model_name') notnull" json:"col_name"`
	Label           string   `json:"label,omitempty"`
	SQLType         SQLType  `xorm:"sql_type json" json:"sql_type,omitempty"`
	IsPrimaryKey    bool     `json:"is_primary_key,omitempty"`
	IsAutoIncrement bool     `json:"is_auto_increment,omitempty"`
	Comment         string   `json:"comment,omitempty"`
	Index           bool     `json:"index,omitempty"`
	Nullable        bool     `json:"nullable,omitempty"`
	Options         []string `json:"options,omitempty"`
	DefaultValue    string   `json:"default,omitempty"`
}

type Model struct {
	ID        uint      `xorm:"id pk autoincr" json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	TableName struct{}  `json:"table_name"`
	Columns   []*Column `json:"columns,omitempty"`
}

func NewColumn(name string, label string) *Column {
	return &Column{
		Name:  name,
		Label: label,
	}
}

func (column *Column) ConvertTo() *schemas.Column {
	newColumn := schemas.NewColumn(column.Name,
		"",
		schemas.SQLType(column.SQLType),
		0, 0, column.Nullable)
	// 注意，NewColumn之后的default
	// newColumn.MapType == schemas.TWOSIDES // 默认值
	// newColumn.Default == true
	newColumn.Comment = column.Comment
	newColumn.IsAutoIncrement = column.IsAutoIncrement
	newColumn.IsPrimaryKey = column.IsPrimaryKey

	return newColumn
}

func (column *Column) ConvertFrom(xColumn *schemas.Column) *Column {
	column.Name = xColumn.Name
	column.Nullable = xColumn.Nullable
	column.Comment = xColumn.Comment
	column.IsAutoIncrement = xColumn.IsAutoIncrement
	column.IsPrimaryKey = xColumn.IsPrimaryKey
	column.SQLType = SQLType(xColumn.SQLType)
	return column
}

func (column *Column) Compare(xColumn *schemas.Column) bool {
	return column.Name == xColumn.Name &&
		column.Nullable == xColumn.Nullable &&
		column.Comment == xColumn.Comment &&
		column.IsAutoIncrement == xColumn.IsAutoIncrement &&
		reflect.DeepEqual(column.DefaultValue, xColumn.Default) &&
		reflect.DeepEqual(column.SQLType, xColumn.SQLType) &&
		column.IsPrimaryKey == xColumn.IsPrimaryKey &&
		column.IsAutoIncrement == xColumn.IsAutoIncrement
	// column.SQLType == SQLType(xColumn.SQLType)
}
func Compare(x, y *schemas.Column) bool {
	return reflect.DeepEqual(x, y)
}
