package dsl

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"xorm.io/xorm/schemas"
)

type tag struct {
	name   string
	params []string
}

type Context struct {
	tag
	tagUname        string
	preTag, nextTag string
	table           *schemas.Table
	col             *schemas.Column
	fieldValue      reflect.Value
	isIndex         bool
	isUnique        bool
	indexNames      map[string]int
	hasCacheTag     bool
	hasNoCacheTag   bool
	ignoreNext      bool
	isUnsigned      bool
}

var tagTypeMapper = map[string]reflect.Type{
	"TEXT":   reflect.TypeOf(string("")),
	"BIGINT": reflect.TypeOf(int64(0)),
	"INT":    reflect.TypeOf(int32(0)),
}

func (ctx *Context) ToTagsColumn() []string {
	col := ctx.col
	tags := make([]string, 0)
	if col.IsJSON {
		tags = append(tags, "JSON")
	}
	if len(col.EnumOptions) > 0 {
		enumItems := make([]string, 0)
		for k, _ := range col.EnumOptions {
			enumItems = append(enumItems, k)
		}
		tags = append(tags, schemas.Enum+"("+strings.Join(enumItems, ", ")+")")
	}
	if len(col.SetOptions) > 0 {
		setItems := make([]string, 0)
		for k, _ := range col.SetOptions {
			setItems = append(setItems, k)
		}
		tags = append(tags, schemas.Set+"("+strings.Join(setItems, ", ")+")")
	}
	// index
	if len(col.Indexes) > 0 {
		indexItems := make([]string, 0)
		for k, _ := range col.Indexes {
			indexItems = append(indexItems, k)
		}
		tags = append(tags, "INDEX("+strings.Join(indexItems, ", ")+")")
	}

	if !col.DefaultIsEmpty {
		if col.Length2 != 0 {
			tags = append(tags, fmt.Sprintf("(%d,%d)", col.Length, col.Length2))
		} else if col.Length != 0 {
			tags = append(tags, fmt.Sprintf("(%d)", col.Length))
		}
	}

	if !col.Nullable { // not null
		if col.IsPrimaryKey {
			tags = append(tags, "PK")
		}
		if col.IsAutoIncrement {
			tags = append(tags, "AUTOINCR")
		}
	}

	if col.IsPrimaryKey && !col.Nullable {
		tags = append(tags, "PK")
	} else if !col.Nullable {
		tags = append(tags, "NOT NULL")
	}

	switch col.TimeZone {
	case time.Local:
		tags = append(tags, "LOCAL")
		break
	case time.UTC:
		tags = append(tags, "UTC")
	}

	// EXTENDS
	if len(col.Comment) > 0 {
		tags = append(tags, fmt.Sprintf("COMMENT '%s'", col.Comment))
	}

	return tags
}
