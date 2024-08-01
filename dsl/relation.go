package dsl

// 表之间关系定义

type Relation struct {
	Type    string
	Model   string
	Key     string
	Foreign string
	Query   string
}
