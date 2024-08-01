package models_test

import (
	"encoding/json"
	"github.com/everpan/mdbg/models"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"xorm.io/xorm/schemas"
)

func TestColumn_ConvertFrom(t *testing.T) {
	type fields struct {
		ID              uint64
		ModelID         uint64
		Name            string
		Label           string
		SQLType         schemas.SQLType
		IsPrimaryKey    bool
		IsAutoIncrement bool
		Comment         string
		Index           bool
		Nullable        bool
		Options         []string
		Default         string
	}
	type args struct {
		xColumn *schemas.Column
	}
	var tests []struct {
		name   string
		fields fields
		args   args
		want   *models.Column
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := &models.Column{
				ID:              tt.fields.ID,
				ModelID:         tt.fields.ModelID,
				Name:            tt.fields.Name,
				Label:           tt.fields.Label,
				SQLType:         tt.fields.SQLType,
				IsPrimaryKey:    tt.fields.IsPrimaryKey,
				IsAutoIncrement: tt.fields.IsAutoIncrement,
				Comment:         tt.fields.Comment,
				Index:           tt.fields.Index,
				Nullable:        tt.fields.Nullable,
				Options:         tt.fields.Options,
				Default:         tt.fields.Default,
			}
			if got := column.ConvertFrom(tt.args.xColumn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColumn_ConvertTo(t *testing.T) {
	testData := []struct {
		name string
		col  *models.Column
		want *schemas.Column
	}{
		{col: &models.Column{
			Name:    "test-name",
			SQLType: schemas.SQLType{Name: "type", DefaultLength: 12, DefaultLength2: 34},
		}, want: func() *schemas.Column {
			c := schemas.NewColumn("test-name", "",
				schemas.SQLType{Name: "type", DefaultLength: 12, DefaultLength2: 34},
				0, 0, false)
			return c
		}()},
	}
	tests := testData
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.col.ConvertTo(); !models.Compare(got, tt.want) {
				t.Errorf("CaseName:%s\tConvertTo() = \n%v,\nwant\n%v", tt.name, got, tt.want)
				j1, _ := json.Marshal(got)
				j2, _ := json.Marshal(tt.want)
				t.Errorf("got:\n %s\nwant:\n %s", string(j1), string(j2))
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNewColumn(t *testing.T) {
	type args struct {
		name  string
		label string
	}
	tests := []struct {
		name string
		args args
		want *models.Column
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := models.NewColumn(tt.args.name, tt.args.label); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}
