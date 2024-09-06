package event

type Event struct {
	EventID         uint64 `json:"event_id"`
	EventType       string `json:"event_type"`
	EntityType      string `json:"entity_type"`
	EntityId        uint64 `json:"entity_id"`
	EventData       string `json:"event_data"`
	TriggeringEvent string `json:"triggering_event"` // 重复事件与消息的检测
}

type Entity struct {
	EntityType    string `json:"entity_type"`
	EntityId      uint64 `json:"entity_id"`
	EntityVersion uint32 `json:"entity_version"`
	// EntitySchema  string `json:"entity_schema"`
}

type Snapshot struct {
	Entity
	SnapshotType     string `json:"snapshot_type"`     // 类型
	SnapshotMarshall string `json:"snapshot_marshall"` // 序列化表示
	TriggeringEvent  string `json:"triggering_event"`
}

type EventInterface interface {
	Add(e *Event) error
	Fetch(pk uint64) *Event
	FetchGte(pk uint64, limit int32) []*Event
}
