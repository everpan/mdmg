package ctx

import (
	"github.com/stretchr/testify/assert"
	v8 "rogchap.com/v8go"
	"testing"
	"xorm.io/xorm"
)

func (c *IcContext) SetEngine(engine *xorm.Engine) {
	c.engine = engine
}
func (c *IcContext) SetV8Ctx(context *v8.Context) {
	c.v8Ctx = context
}

func TestIcPage_CalCountOffset(t *testing.T) {
	type Page struct {
		Size        int
		No          int
		Count       int
		RecordCount int
	}
	tests := []struct {
		name        string
		page        Page
		recordCount int
		wantOffset  int
		pageCount   int
	}{
		{"page 0", Page{10, 0, 0, 0}, 20, 0, 2},
		{"page count 3", Page{10, 0, 0, 0}, 21, 0, 3},
		{"page count 3", Page{10, 0, 0, 0}, 29, 0, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &IcPage{
				PageSize:    tt.page.Size,
				PageNo:      tt.page.No,
				PageCount:   tt.page.Count,
				RecordCount: tt.page.RecordCount,
			}
			if gotOffset := p.CalCountOffset(tt.recordCount); gotOffset != tt.wantOffset {
				t.Errorf("CalCountOffset() = %v, want %v", gotOffset, tt.wantOffset)
			}
			assert.Equal(t, tt.pageCount, p.PageCount)
			assert.Equal(t, tt.recordCount, p.RecordCount)
		})
	}
}
