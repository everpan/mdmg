package models_test

import (
	"encoding/json"
	"github.com/everpan/mdmg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"xorm.io/xorm/schemas"
)

func TestColumn_ConvertTo(t *testing.T) {
	testData := []struct {
		name string
		col  *models.Column
		want *schemas.Column
	}{
		{col: &models.Column{
			Name:    "test-name",
			SQLType: (models.SQLType)(schemas.SQLType{Name: "type", DefaultLength: 12, DefaultLength2: 34}),
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
