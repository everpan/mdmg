package dsl

import (
	"encoding/json"
	"strings"
	"xorm.io/xorm/schemas"
)

type mySQLType struct {
	Name           string `json:"name,omitempty"`
	DefaultLength  int64  `json:"defaultLength,omitempty"`
	DefaultLength2 int64  `json:"defaultLength2,omitempty"`
}

func ColumnMarshal(col *schemas.Column) string {
	sb := strings.Builder{}
	sb.WriteString("{\"name\":\"")
	sb.WriteString(col.Name)
	sb.WriteString("\",\"sqlType\":")
	d, _ := json.Marshal(mySQLType(col.SQLType))
	sb.Write(d)

	if col.IsPrimaryKey {
		sb.WriteString(",\"primary\": true")
	}
	if !col.Nullable {
		sb.WriteString(",\"nullable\": false")
	}
	sb.WriteRune('}')
	return sb.String()
}

func ColumnUnmarshal(json string) *schemas.Column {
	return nil
}
