package dsl

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"xorm.io/xorm/schemas"
)

func TestColumnMarshal(t *testing.T) {
	tests := []struct {
		name string
		col  *schemas.Column
		want string
	}{
		{"default column",
			schemas.NewColumn("col_name", "field_name",
				schemas.SQLType{Name: "int", DefaultLength: 1, DefaultLength2: 2}, 0, 0, false),
			`{"name":"col_name","sqlType":{"name":"int","defaultLength":1,"defaultLength2":2},"nullable": false}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ColumnMarshal(tt.col)
			assert.Equal(t, tt.want, got)
		})
	}
}
