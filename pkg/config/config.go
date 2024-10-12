package config

import (
	"encoding/json"
	"fmt"
	"github.com/everpan/mdmg/pkg/config/values"
	"strings"
)

type ItemDesc struct {
	Item string `json:"-"` // the key
	Desc string `json:"desc"`
}

type EnumDesc map[string]string

type Schema struct {
	Item     string        `json:"-"` // the key
	Desc     string        `json:"desc"`
	Type     values.VType  `json:"type"`                // category string number enum
	EnumDesc EnumDesc      `json:"enum-desc,omitempty"` // if type is enum , give each desc of enum value, item is the enum value.
	Value    values.IValue `json:"-"`
	// Default  values.IValue `json:"default,omitempty"`
	Default any `json:"default,omitempty"`
}

func (sc *Schema) GetValue() any {
	//if sc.Value == nil {
	//	if sc.Default != nil {
	//		return sc.Default
	//	}
	//} else {
	//	return sc.Value.Value()
	//}
	if sc.Value != nil {
		return sc.Value.Value()
	}
	return sc.Default
}

func (sc *Schema) Update(schema *Schema, forceUpdateValue bool) {
	sc.Desc = schema.Desc
	sc.Type = schema.Type
	sc.EnumDesc = schema.EnumDesc
	if forceUpdateValue {
		sc.Value = schema.Value
	}
	sc.Default = schema.Default
}

func (sc *Schema) SetValue(val any) error {
	v, e := values.CreateValue(sc.Type, fmt.Sprintf("%v", val))
	if e != nil {
		return e
	}
	sc.Value = v
	return nil
}

func (sc *Schema) ClearValue() {
	sc.Value = nil
}

func NewIConfig() *IcConfig {
	return &IcConfig{
		confMap: make(map[string]*IcSectionConfig),
	}
}

func NewItemSchema(item, desc string, typ values.VType, val, defVal string) *Schema {
	var (
		// err error
		sVal, dVal values.IValue
		dv         any
	)
	if val != "" {
		sVal, _ = values.CreateValue(typ, val)
	}
	dVal, _ = values.CreateValue(typ, defVal)
	if dVal != nil {
		dv = dVal.Value()
	}
	return &Schema{strings.TrimSpace(item),
		desc, typ, nil, sVal, dv}
}

func init() {
	GlobalConfig.NewSection("system", "系统")
}

func (c *IcConfig) NewSection(section, desc string) *IcSectionConfig {
	seConf, ok := c.confMap[section]
	if !ok {
		seConf = &IcSectionConfig{description: desc, sMap: make(map[string]*Schema)}
		c.confMap[section] = seConf
	}
	return seConf
}

func (c *IcConfig) Section(sec string) *IcSectionConfig {
	seConf, ok := c.confMap[sec]
	if !ok {
		return nil
	}
	return seConf
}

func SplitSectionKey(key string) (section, subKey string) {
	key1 := strings.TrimSpace(key)
	dotPos := strings.IndexByte(key1, '.')
	if dotPos > -1 {
		section = key1[:dotPos]
		subKey = key1[dotPos+1:]
		return
	}
	return key, ""
}

func (c *IcConfig) GetValue(key string) any {
	sec, subKey := SplitSectionKey(key)
	if subKey == "" {
		return nil
	}
	secConf := c.Section(sec)
	if secConf == nil {
		return nil
	}
	return secConf.GetValue(subKey)
}

func (c *IcConfig) SetValue(key string, val any) error {
	sec, subKey := SplitSectionKey(key)
	if subKey == "" {
		return fmt.Errorf("invalid section key: %s", key)
	}
	sc := c.Section(sec)
	if sc == nil {
		return fmt.Errorf("invalid schema of section: %s", sec)
	}
	return sc.SetValue(subKey, val)
}

func (c *IcConfig) OutputMap() any {
	m := make(map[string]any)
	for k, v := range c.confMap {
		m[k] = v.OutputMap()
	}
	return m
}

func (c *IcConfig) OutputValues(flat bool) []byte {
	if flat {
		return c.outputFlatValues()
	}
	m := c.OutputMap()
	d, _ := json.MarshalIndent(m, "", "  ")
	return d
}

// OutputFlatValues 将所有的配置项扁平化为 "key"=value 的形式
func (c *IcConfig) outputFlatValues() []byte {
	m := make(map[string]any)
	for k, v := range c.confMap {
		m2 := v.OutputMap()
		for k2, v2 := range m2 {
			m[k+"."+k2] = v2
		}
	}
	d, _ := json.MarshalIndent(m, "", "  ")
	return d
}
