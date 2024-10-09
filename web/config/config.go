package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type Config struct {
	JSModuleRootPath string
	// module-version/${JSModuleBeckEndDir}/script.js
	JSModuleBeckEndDir string
	HiddenModelVersion bool `yaml:"hidden-model-version"`
}

type ItemDesc struct {
	Item string `json:"-"` // the key
	Desc string `json:"desc"`
}

type EnumDesc map[string]string

func (itemDesc *ItemDesc) Encode() []byte {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(itemDesc.Item)
	buf.Truncate(buf.Len() - 1)
	buf.WriteString(": ")
	encoder.Encode(itemDesc.Desc)
	buf.Truncate(buf.Len() - 1)
	return buf.Bytes()
}

type Schema struct {
	Item     string   `json:"-"` // the key
	Desc     string   `json:"desc"`
	Type     string   `json:"type"`                // category string number enum
	EnumDesc EnumDesc `json:"enum-desc,omitempty"` // if type is enum , give each desc of enum value, item is the enum value.
	Value    string   `json:"value,omitempty"`
	Default  string   `json:"default,omitempty"`
}

type SchemaMap = map[string]*Schema

type IcConfig struct {
	sMap SchemaMap
}

var (
	ICodeGlobalConfig = NewIConfig()
)

type ValueType = string

const (
	CategoryType ValueType = "category"
	StringType   ValueType = "string"
	NumberType   ValueType = "number"
	FloatType    ValueType = "float"
	BooleanType  ValueType = "boolean"
	EnumType     ValueType = "enum"
)

func (sc *Schema) GetValue() string {
	if len(sc.Value) == 0 {
		return sc.Default
	}
	return sc.Value
}
func NewIConfig() *IcConfig {
	return &IcConfig{
		sMap: SchemaMap{},
	}
}
func NewItemSchema(item, desc string, typ ValueType, val, def string) *Schema {
	return &Schema{
		item, desc, typ, nil, val, def,
	}
}

func init() {
	ICodeGlobalConfig.AddSection("system", "系统")
	ICodeGlobalConfig.AddSection("js-module", "JS模块")
	ICodeGlobalConfig.AddStringSchema("js-module.root-path", "根目录", "web/js")
	ICodeGlobalConfig.AddEnumSchema("js-module.version-in-path",
		"模块路径中是否有版本，如：mod-1.3.4", "true",
		EnumDesc{"true": "包含", "false": "不包含，需要另外的映射"})
}
func (c *IcConfig) GetSchema(key string) *Schema {
	sc, ok := c.sMap[key]
	if !ok {
		return nil
	}
	return sc
}
func (c *IcConfig) GetValue(key string) (any, bool) {
	conf, ok := c.sMap[key]
	if ok {
		return conf.GetValue(), ok
	}
	return "", false
}

// OutputSchema 输出schema信息
func (c *IcConfig) OutputSchema() []byte {
	if len(c.sMap) == 0 {
		return []byte("{}")
	}
	buf := bytes.NewBufferString("{\n")
	for k, c := range c.sMap {
		buf.WriteString(" \"")
		buf.WriteString(k)
		buf.WriteString("\": ")
		data, _ := json.Marshal(c)
		buf.Write(data)
		buf.WriteString(",\n")
	}
	if len(c.sMap) > 0 {
		buf.Truncate(buf.Len() - 2)
	}
	buf.WriteString("\n}")
	return buf.Bytes()
}

// ParseSchema 解析schema，形成默认配置
func (c *IcConfig) ParseSchema(data []byte) error {
	return json.Unmarshal(data, &c.sMap)
}

func (c *IcConfig) LoadKeyValue(data []byte) error {
	return nil
}

func (c *IcConfig) Validate(data []byte) []error {
	kv := make(map[string]any)
	if err := json.Unmarshal(data, &kv); err != nil {
		return []error{err}
	}
	errs := make([]error, 0)
	for k, v := range kv {
		sc := c.GetSchema(k)
		if sc == nil {
			errs = append(errs, fmt.Errorf("not found schema of `%v`", k))
		}
		// check v and sc.Type is same
		switch v.(type) {
		case string:
			if sc.Type != StringType {
				errs = append(errs, fmt.Errorf("not support type of `%v`, need string", k))
			}
		case int64, int32:
			if sc.Type != NumberType {
				errs = append(errs, fmt.Errorf("not support type of `%v`, need number", k))
			}
		case float64, float32:
			if sc.Type != FloatType {
				errs = append(errs, fmt.Errorf("not support type of `%v`, need float", k))
			}
		case bool:
			if sc.Type != BooleanType {
				errs = append(errs, fmt.Errorf("not support type of `%v`, need boolean", k))
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func FetchValueByType(val, typ string) any {
	if typ == CategoryType {
		return nil
	}
	switch typ {
	case NumberType:
		v, _ := strconv.ParseInt(val, 10, 64)
		return v
	case FloatType:
		v, _ := strconv.ParseFloat(val, 64)
		return v
	default:
		return val
	}
}

func (c *IcConfig) OutputMap() map[string]any {
	kv := make(map[string]any)
	for _, item := range c.sMap {
		if item.Type == CategoryType {
			continue
		}
		strVal := item.Value
		if len(item.Value) == 0 {
			strVal = item.Default
		}
		kv[item.Item] = FetchValueByType(strVal, item.Type)
	}
	return kv
}

// OutputKeyValue 以kv的形式输出配置项; section is ignore
func (c *IcConfig) OutputKeyValue() []byte {
	kv := c.OutputMap()
	data, _ := json.MarshalIndent(kv, "", " ")
	return data
}

func (c *IcConfig) AddSchema(schema *Schema) {
	c.sMap[schema.Item] = schema
}
func (c *IcConfig) AddSection(sec, desc string) {
	c.sMap[sec] = NewItemSchema(sec, desc, CategoryType, "", "")
}
func (c *IcConfig) AddStringSchema(item, desc, defaultValue string) {
	c.sMap[item] = NewItemSchema(item, desc, StringType, "", defaultValue)
}
func (c *IcConfig) AddNumberSchema(item, desc, defaultValue string) {
	c.sMap[item] = NewItemSchema(item, desc, NumberType, "0", defaultValue)
}
func (c *IcConfig) AddFloatSchema(item, desc, defaultValue string) {
	c.sMap[item] = NewItemSchema(item, desc, FloatType, "0.0", defaultValue)
}
func (c *IcConfig) AddBooleanSchema(item, desc, defaultValue string) {
	c.sMap[item] = NewItemSchema(item, desc, BooleanType, "false", defaultValue)
}

func (c *IcConfig) AddEnumSchema(item, desc, defaultValue string, enumDesc EnumDesc) {
	schema := NewItemSchema(item, desc, EnumType, "", defaultValue)
	schema.EnumDesc = enumDesc
	c.sMap[item] = schema
}

var DefaultConfig = Config{
	JSModuleRootPath:   "web/script_module",
	JSModuleBeckEndDir: "backend",
}
