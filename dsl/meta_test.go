package dsl

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"xorm.io/xorm"
	"xorm.io/xorm/core"
	"xorm.io/xorm/dialects"
	"xorm.io/xorm/schemas"
)

var (
	basePath = "examples"
)

func TestReadMetaJson(t *testing.T) {
	data := []struct {
		Filename string
		want     []string
	}{
		{
			Filename: "examples/product.json",
			want: []string{
				`{"name":"user_id","label":"用户ID","type":"integer","index":true}`,
				`{"name":"name","label":"名称","type":"varchar","length1":128,"index":true}`,
			},
		},
	}
	for _, td := range data {
		fdata, err := os.ReadFile(td.Filename)
		if err != nil {
			t.Fatal(err)
		}
		var meta = Meta{}
		err = json.Unmarshal(fdata, &meta)
		if err != nil {
			t.Fatal(err)
		}
		data2, _ := json.Marshal(meta)
		metaJson := string(data2)
		for _, want := range td.want {
			assert.Contains(t, metaJson, want, metaJson)
		}
	}
}

func TestCreateTableFromMeta(t *testing.T) {
	mData, err := os.ReadFile("examples/product.json")
	if err != nil {
		t.Fatal(err)
	}
	meta := Meta{}
	err = json.Unmarshal(mData, &meta)
	if err != nil {
		t.Fatal(err)
	}
	table := schemas.NewEmptyTable()
	table.Name = meta.Table.Name
	table.Comment = meta.Table.Comment

	orm, err := xorm.NewEngine("mysql", "root:@tcp(127.0.0.1:3306)/wiz_hr2?charset=utf8")
	if err != nil {
		t.Fatal(err)
	}
	dialect := orm.Dialect()
	sqlStr, err := meta.CreateTableSQL(orm.DB(), dialect)
	assert.NotEmptyf(t, sqlStr, "sql empty!")
}

func fatalErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestExportDbColumnsIntoExamples(t *testing.T) {
	t.Skip("manual run")
	dbName := "wiz_hr2"
	connStr := fmt.Sprintf("root:123456@tcp(127.0.0.1:3306)/%s?charset=utf8", dbName)
	orm, err := xorm.NewEngine("mysql", connStr)
	fatalErr(t, err)
	ms, err := orm.DBMetas()
	fatalErr(t, err)
	dbBasePath := basePath + "/xorm_columns/" + dbName
	err = os.MkdirAll(basePath, os.ModePerm)
	fatalErr(t, err)
	for _, m := range ms {
		for _, c := range m.Columns() {
			fName := fmt.Sprintf("%s/%s.%s.json", dbBasePath, m.Name, c.Name)
			f, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, os.ModePerm)
			defer f.Close()
			fatalErr(t, err)
			d, err := json.MarshalIndent(c, "", " ")
			fatalErr(t, err)
			f.Write(d)
		}
	}
}

func TestColumnsConvertDslToXorm(t *testing.T) {
	tData := []struct {
		name     string
		fileName string
		contains []string
	}{
		{name: "unsigned int", fileName: "dsl_columns/int.json",
			contains: []string{"UINT32", "BIGINT(20) UNSIGNED", "INT UNSIGNED"}},
	}
	db, err := core.Open("mysql", "")
	fatalErr(t, err)
	dialect, err := dialects.OpenDialect("mysql", "")
	fatalErr(t, err)
	for _, td := range tData {
		t.Run(td.name, func(t *testing.T) {
			meta := Meta{} //.Columns
			fileName := fmt.Sprintf("%s/%s", basePath, td.fileName)
			fData, err := os.ReadFile(fileName)
			fatalErr(t, err)
			err = json.Unmarshal(fData, &meta)
			fatalErr(t, err)
			cols := meta.Columns
			assert.Equal(t, len(cols), len(td.contains),
				"contains length must equal columns length")
			for i, col := range cols {
				xCol := col.ConvertXormColumn(db)
				cStr, err := dialects.ColumnString(dialect, xCol, true, false)
				fatalErr(t, err)
				assert.Contains(t, cStr, td.contains[i], "%s must contain %s", cStr, td.contains[i])
			}
		})
	}
}

func TestXormColumnConvertTo(t *testing.T) {
	tData := []struct {
		name     string
		fileName string
		want     string
	}{
		{name: "unsigned int", fileName: "wiz_hr2/_prisma_migrations.applied_steps_count.json", want: "unsigned int"},
		{name: "date time", fileName: "wiz_hr2/_prisma_migrations.started_at.json", want: "datetime"},
		{name: "varchar 5000", fileName: "wiz_hr2/hr_workflow.path.json", want: "varchar"},
		{name: "varchar common", fileName: "microi_empty/diy_component.FieldType.json", want: "varchar"},
		{name: "varchar primary key", fileName: "microi_empty/diy_component.Id.json", want: "varchar"},
	}
	for _, td := range tData {
		t.Run(td.name, func(t *testing.T) {
			fileName := fmt.Sprintf("%s/xorm_columns/%s", basePath, td.fileName)
			fData, err := os.ReadFile(fileName)
			fatalErr(t, err)
			xCol := XormColumn{}
			err = json.Unmarshal(fData, &xCol)
			fatalErr(t, err)
			col := xCol.ConvertDslColumn()
			assert.Equal(t, xCol.SQLType.Name, strings.ToUpper(col.Type))
			assert.Equal(t, xCol.Length, col.Length1)
			assert.Equal(t, xCol.IsPrimaryKey, col.IsPrimaryKey)
			assert.Equal(t, xCol.IsAutoIncrement, col.IsAutoIncrement)
			assert.Equal(t, col.Type, td.want)
			// t.Log(col)
		})
	}
}

func TestColumn_ConvertXormColumn(t *testing.T) {
	t.Skip("not implemented yet")
	tests := []struct {
		name   string
		column *Column
		db     *core.DB
		want   *schemas.Column
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.column.ConvertXormColumn(tt.db), "%s ConvertXormColumn()", tt.name)
		})
	}
}
