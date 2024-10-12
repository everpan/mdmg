package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/everpan/mdmg/pkg/config/values"
)

type IcSectionConfig struct {
	description string
	sMap        map[string]*Schema
}

type IcConfig struct {
	confMap map[string]*IcSectionConfig
}

func (c *IcSectionConfig) GetSchema(key string) *Schema {
	sc, ok := c.sMap[key]
	if !ok {
		return nil
	}
	return sc
}

func (c *IcSectionConfig) GetValue(key string) any {
	schema, ok := c.sMap[key]
	if ok {
		return schema.GetValue()
	}
	return nil
}

// ExportSchema 输出schema信息
func (c *IcSectionConfig) ExportSchema(flat bool) []byte {
	if len(c.sMap) == 0 {
		return []byte("{}")
	}
	buf := bytes.NewBufferString("{\n")
	for k, schema := range c.sMap {
		buf.WriteString(" \"")
		buf.WriteString(k)
		buf.WriteString("\": ")
		data, _ := json.Marshal(schema)
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
func (c *IcSectionConfig) ParseSchema(data []byte) error {
	return json.Unmarshal(data, &c.sMap)
}

func (c *IcSectionConfig) LoadKeyValue(data []byte) error {
	return nil
}

func (c *IcSectionConfig) OutputMap() map[string]any {
	kv := make(map[string]any)
	for _, item := range c.sMap {
		kv[item.Item] = item.GetValue()
	}
	return kv
}

// OutputKeyValue 以kv的形式输出配置项; description is ignore
func (c *IcSectionConfig) OutputKeyValue() []byte {
	kv := c.OutputMap()
	data, _ := json.MarshalIndent(kv, "", " ")
	return data
}

func (c *IcSectionConfig) AddSchema(schema *Schema) {
	c.sMap[schema.Item] = schema
}
func (c *IcSectionConfig) ReplaceSchema(schema *Schema) {
	sc, ok := c.sMap[schema.Item]
	if !ok {
		c.AddSchema(schema)
		return
	}
	sc.Update(schema, false)
}
func (c *IcSectionConfig) ReplaceSchemaAndValue(schema *Schema) {
	sc, ok := c.sMap[schema.Item]
	if !ok {
		c.AddSchema(schema)
		return
	}
	sc.Update(schema, true)
}
func (c *IcSectionConfig) AddStringSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(item, desc, values.VStringT, "", defaultValue))
}
func (c *IcSectionConfig) AddNumberSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(item, desc, values.VNumberT, "", defaultValue))
}
func (c *IcSectionConfig) AddFloatSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(item, desc, values.VFloatT, "", defaultValue))
}
func (c *IcSectionConfig) AddBooleanSchema(item, desc, defaultValue string) {
	c.ReplaceSchema(NewItemSchema(item, desc, values.VBooleanT, "", defaultValue))
}

func (c *IcSectionConfig) AddEnumSchema(item, desc string, valType values.VType, defaultValue string, enumDesc EnumDesc) {
	schema := NewItemSchema(item, desc, values.CompositeT(values.VEnumT, valType), "", defaultValue)
	schema.EnumDesc = enumDesc
	c.ReplaceSchema(schema)
}

func (c *IcSectionConfig) SetValue(key string, val any) error {
	sc := c.GetSchema(key)
	if sc == nil {
		return fmt.Errorf("no schema for key: %s", key)
	}
	return sc.SetValue(val)
}
