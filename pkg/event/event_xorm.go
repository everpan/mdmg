package event

import (
	"strconv"
	"xorm.io/xorm"
)

const (
	EventTable    = "icode_event"
	EntityTable   = "icode_entity"
	SnapshotTable = "icode_snapshot"
)

type XORMEvent struct {
	engine *xorm.Engine
}

func NewXORMEvent() *XORMEvent {
	return &XORMEvent{}
}

func NewXORMEventWithEngine(engine *xorm.Engine) *XORMEvent {
	x := NewXORMEvent()
	x.SetEngine(engine)
	x.setup()
	return x
}

func (x *XORMEvent) setup() {
	x.engine.Table(EventTable).CreateTable(&Event{})
	x.engine.Table(EntityTable).CreateTable(&Entity{})
	x.engine.Table(SnapshotTable).CreateTable(&Snapshot{})
}

func (x *XORMEvent) SetEngine(e *xorm.Engine) {
	x.engine = e
}

func (x *XORMEvent) Driver() string {
	return "xorm"
}

func (x *XORMEvent) Add(e *Event) error {
	e.EventId = *new(uint64)
	_, err := x.engine.Table("icode_event").Insert(e)
	if nil != err {
		x.engine.Logger().Errorf("Add event: %v fail : %v", e, err)
		return err
	}
	Pub(e)
	return nil
}

func (x *XORMEvent) MaxId() uint64 {
	sql := "SELECT MAX(event_id) as m FROM icode_event"
	r, err := x.engine.Query(sql)
	if err != nil {
		x.engine.Logger().Errorf("Get max id error: %v", err)
		return 0
	}
	m, _ := strconv.ParseUint(string(r[0]["m"]), 10, 64)
	return m
}

func (x *XORMEvent) Fetch(eventId uint64) *Event {
	e := &Event{EventId: eventId}
	b, err := x.engine.Table(EventTable).Get(e)
	if nil != err {
		x.engine.Logger().Errorf("Fetch event: %v fail : %v", e, err)
		return nil
	}
	if b {
		return e
	}
	return nil
}

func (x *XORMEvent) FetchGte(eventId uint64, limit int32) []*Event {
	if int32(0) == limit {
		limit = 20
	}
	return nil
}
