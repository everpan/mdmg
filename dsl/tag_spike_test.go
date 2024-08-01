package dsl

import (
	"testing"
	"xorm.io/xorm/caches"
	"xorm.io/xorm/dialects"
	"xorm.io/xorm/names"
	"xorm.io/xorm/tags"
)

/*
func parseFieldWithTags(parser *xormTags.Parser,

		table *schemas.Table, fieldIndex int, field reflect.StructField,
		fieldValue reflect.Value, tags []tag) (*schemas.Column, error) {
		col := &schemas.Column{
			FieldName:       field.Name,
			FieldIndex:      []int{fieldIndex},
			Nullable:        true,
			IsPrimaryKey:    false,
			IsAutoIncrement: false,
			MapType:         schemas.TWOSIDES,
			Indexes:         make(map[string]int),
			DefaultIsEmpty:  true,
		}

		ctx := xormTags.Context{
			table:      table,
			col:        col,
			fieldValue: fieldValue,
			indexNames: make(map[string]int),
			parser:     parser,
		}

		for j, tag := range tags {
			if ctx.ignoreNext {
				ctx.ignoreNext = false
				continue
			}

			ctx.tag = tag
			ctx.tagUname = strings.ToUpper(tag.name)

			if j > 0 {
				ctx.preTag = strings.ToUpper(tags[j-1].name)
			}
			if j < len(tags)-1 {
				ctx.nextTag = tags[j+1].name
			} else {
				ctx.nextTag = ""
			}

			if h, ok := parser.handlers[ctx.tagUname]; ok {
				if err := h(&ctx); err != nil {
					return nil, err
				}
			} else {
				if strings.HasPrefix(ctx.tag.name, "'") && strings.HasSuffix(ctx.tag.name, "'") {
					col.Name = ctx.tag.name[1 : len(ctx.tag.name)-1]
				} else {
					col.Name = ctx.tag.name
				}
			}

			if ctx.hasCacheTag {
				if parser.cacherMgr.GetDefaultCacher() != nil {
					parser.cacherMgr.SetCacher(table.Name, parser.cacherMgr.GetDefaultCacher())
				} else {
					parser.cacherMgr.SetCacher(table.Name, caches.NewLRUCacher2(caches.NewMemoryStore(), time.Hour, 10000))
				}
			}
			if ctx.hasNoCacheTag {
				parser.cacherMgr.SetCacher(table.Name, nil)
			}
		}

		if col.SQLType.Name == "" {
			var err error
			col.SQLType, err = parser.getSQLTypeByType(field.Type)
			if err != nil {
				return nil, err
			}
		}
		if ctx.isUnsigned && col.SQLType.IsNumeric() && !strings.HasPrefix(col.SQLType.Name, "UNSIGNED") {
			col.SQLType.Name = "UNSIGNED " + col.SQLType.Name
		}

		parser.dialect.SQLType(col)
		if col.Length == 0 {
			col.Length = col.SQLType.DefaultLength
		}
		if col.Length2 == 0 {
			col.Length2 = col.SQLType.DefaultLength2
		}
		if col.Name == "" {
			col.Name = parser.columnMapper.Obj2Table(field.Name)
		}

		if ctx.isUnique {
			ctx.indexNames[col.Name] = schemas.UniqueType
		} else if ctx.isIndex {
			ctx.indexNames[col.Name] = schemas.IndexType
		}

		for indexName, indexType := range ctx.indexNames {
			addIndex(indexName, table, col, indexType)
		}

		return col, nil
	}
*/

func TestTagStrToColumn(t *testing.T) {
	_ = tags.NewParser(
		"xorm",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.SnakeMapper{},
		caches.NewManager(),
	)
	// parser.ParseFieldWithTags(nil,1,)
}
