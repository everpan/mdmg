package event

type Event struct {
	EventId         uint64 `json:"event_id" xorm:"pk autoincr"`
	EventType       string `json:"event_type"`
	EntityType      string `json:"entity_type"`
	EntityId        uint64 `json:"entity_id"`
	EventData       string `json:"event_data"`
	CreatedTime     int64  `json:"created_time" xorm:"created"`
	TriggeringEvent string `json:"triggering_event"` // 重复事件与消息的检测
}

type Entity struct {
	EntityTypeId  uint64 `json:"entity_type_id" xorm:"pk autoincr"` // 实体类别
	EntityType    string `json:"entity_type"`
	EntityVersion uint32 `json:"entity_version"`
	// EntitySchema  string `json:"entity_schema"`
}

type Snapshot struct {
	Entity
	SnapshotType     string `json:"snapshot_type"`     // 类型
	SnapshotMarshall string `json:"snapshot_marshall"` // 序列化表示
	TriggeringEvent  string `json:"triggering_event"`
}

type IEvent interface {
	Driver() string
	Add(e *Event) error
	Fetch(eventId uint64) *Event
	FetchGte(eventId uint64, limit int32) []*Event
}
