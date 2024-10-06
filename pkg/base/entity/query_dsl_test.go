package entity

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueryDMLInitJustIncludeSelect(t *testing.T) {
	query := QueryDML{}
	d, _ := json.Marshal(query)
	// t.Log(string(data))
	assert.Equal(t, string(d), `{"select":null}`)
}

func TestQueryDSL_BuildWhere(t *testing.T) {
	ws := []*WhereDML{
		{"col1", "val1", "", "", nil, 0},
		{"col2", "val2", "", "", nil, 0},
	}
	w0 := &WhereDML{"col0", "val0", "", "", ws, 0}
	wor := &WhereDML{"col_or", "val_or", "", "or", nil, 0}
	tests := []struct {
		name   string
		wheres WheresDML
		result string
	}{
		{"signal condition", ws[0:1], `where col1 = "val1"`},
		{"and two condition", ws[0:2], `where col1 = "val1" and col2 = "val2"`},
		{"and or condition", []*WhereDML{ws[0], wor, ws[1]}, `where col1 = "val1" or col_or = "val_or" and col2 = "val2"`},
		{"sub condition", []*WhereDML{ws[0], w0, ws[1]},
			`where col1 = "val1" and (col1 = "val1" and col2 = "val2") and col2 = "val2"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewBuilder()
			buf := bytes.Buffer{}
			q.Clear()
			q.whereSQL(&buf, tt.wheres)
			assert.Equalf(t, tt.result, buf.String(), "Where()")
		})
	}
}

func TestBuilder_selectSQL(t *testing.T) {
	tests := []struct {
		name     string
		selItems SelectDML
		want     string
	}{
		{"0 item", nil, "select "},
		{"1 item", []string{"a"}, "select a"},
		{"3 items", []string{"a", "b", "c"}, "select a,b,c"},
		{"alias", []string{"a as a"}, "select a as a"},
		{"table scope item", []string{"t.a as a"}, "select t.a as a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewBuilder()
			buf := bytes.Buffer{}
			q.selectSimpleSQL(&buf, tt.selItems)
			assert.Equalf(t, tt.want, buf.String(), "selectSimpleSQL(%v)", tt.selItems)
		})
	}
}
