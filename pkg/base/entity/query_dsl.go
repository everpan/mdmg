package entity

import (
	"bytes"
	"strconv"
)

/** 将json查询结构转换为sql */

type WhereDML struct {
	Col     string    `json:"col"`
	Val     string    `json:"val" json:"val,omitempty"`
	Op      string    `json:"op,omitempty"`
	Combine string    `json:"combine,omitempty"`
	Wheres  WheresDML `json:"wheres,omitempty"`
	ValType int       `json:"val_type,omitempty"` // default 0:string, 1:int
}
type WheresDML []*WhereDML
type SelectDML []string
type OrderDML struct {
	Col    string `json:"col"`
	Option string `json:"opt"`
}
type OrdersDML []*OrderDML
type LimitDML []uint32
type QueryDML struct {
	Select SelectDML `json:"select"`
	Wheres WheresDML `json:"wheres,omitempty"`
	Orders OrdersDML `json:"orders,omitempty"`
	Limit  LimitDML  `json:"limit,omitempty"`
}
type Builder struct {
	buf        *bytes.Buffer
	query      QueryDML
	isSubWhere bool
}

func NewBuilder() *Builder {
	return &Builder{
		buf: &bytes.Buffer{},
		query: QueryDML{
			Select: make(SelectDML, 0),
			Wheres: make(WheresDML, 0),
			Orders: make(OrdersDML, 0),
			Limit:  make(LimitDML, 0),
		},
		isSubWhere: false,
	}
}
func (builder *Builder) Where(wheres WheresDML) *Builder {
	builder.query.Wheres = append(builder.query.Wheres, wheres...)
	return builder
}

func (builder *Builder) whereSQL(buf *bytes.Buffer, wheres WheresDML) *Builder {
	if !builder.isSubWhere {
		buf.WriteString("where ")
	}
	builder.isSubWhere = true
	for i, where := range wheres {
		if i > 0 && len(where.Combine) == 0 {
			buf.WriteString(" and ")
		}
		if len(where.Wheres) > 0 {
			buf.WriteByte('(')
			builder.whereSQL(buf, where.Wheres)
			buf.WriteByte(')')
			// buf.WriteByte(' ')
			continue
		}
		if len(where.Combine) > 0 {
			buf.WriteByte(' ')
			buf.WriteString(where.Combine)
			buf.WriteByte(' ')
		}
		buf.WriteString(where.Col)

		if len(where.Op) > 0 {
			buf.WriteByte(' ')
			buf.WriteString(where.Op)
			buf.WriteByte(' ')
		} else {
			buf.WriteString(" = ")
		}
		if where.ValType == 0 {
			buf.WriteByte('"')
			buf.WriteString(where.Val)
			buf.WriteByte('"')
		} else {
			buf.WriteString(where.Val)
		}
	}
	return builder
}
func (builder *Builder) Select(selItems SelectDML) *Builder {
	builder.query.Select = append(builder.query.Select, selItems...)
	return builder
}

func (builder *Builder) selectSQL(buf *bytes.Buffer, selItems SelectDML) *Builder {
	buf.WriteString("select ")
	for i, item := range selItems {
		if i > 0 && len(item) > 0 {
			buf.WriteByte(',')
			buf.WriteString(item)
		} else {
			buf.WriteString(item)
		}
	}
	return builder
}

func (builder *Builder) OrderBy(orderItems OrdersDML) *Builder {
	builder.query.Orders = append(builder.query.Orders, orderItems...)
	return builder
}

func (builder *Builder) orderBySQL(buf *bytes.Buffer, orderItems OrdersDML) *Builder {
	if len(orderItems) > 0 {
		buf.WriteString(" order by ")
	}
	for i, item := range orderItems {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(item.Col)
		buf.WriteByte(' ')
		if len(item.Option) > 0 {
			buf.WriteString(item.Option)
		} else {
			buf.WriteString("asc")
		}
	}
	return builder
}

func (builder *Builder) Limit(limit LimitDML) *Builder {
	builder.query.Limit = append(builder.query.Limit, limit...)
	return builder
}

func (builder *Builder) limitSQL(buf *bytes.Buffer, limit LimitDML) *Builder {
	if len(limit) == 0 {
		return builder
	}
	buf.WriteString(" limit ")
	buf.WriteString(strconv.FormatUint(uint64(limit[0]), 10))
	if len(limit) > 1 {
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatUint(uint64(limit[1]), 10))
	}
	return builder
}

func (builder *Builder) SQL() (string, error) {
	builder.selectSQL(builder.buf, builder.query.Select)
	builder.isSubWhere = false
	builder.whereSQL(builder.buf, builder.query.Wheres)
	builder.orderBySQL(builder.buf, builder.query.Orders)
	builder.limitSQL(builder.buf, builder.query.Limit)
	return builder.buf.String(), nil
}

func (builder *Builder) Clear() {
	builder.isSubWhere = false
	builder.buf.Reset()
	builder.query = QueryDML{
		Select: builder.query.Select[:0],
		Wheres: builder.query.Wheres[:0],
		Orders: builder.query.Orders[:0],
		Limit:  builder.query.Limit[:0],
	}
}
