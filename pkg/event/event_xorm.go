package event

import (
	"xorm.io/xorm"
)

const (
	EventTable    = "icode_event"
	EntityTable   = "icode_entity"
	SnapshotTable = "icode_snapshot"
)

type XORM struct {
	engine *xorm.Engine
}

func NewXORM() *XORM {
	return &XORM{}
}

func NewXORMWithEngine(engine *xorm.Engine) *XORM {
	x := NewXORM()
	x.SetEngine(engine)
	x.setup()
	return x
}

func (x *XORM) setup() {
	x.engine.Table(EventTable).CreateTable(&Event{})
	x.engine.Table(EntityTable).CreateTable(&Entity{})
	x.engine.Table(SnapshotTable).CreateTable(&Snapshot{})
}

func (x *XORM) SetEngine(e *xorm.Engine) {
	x.engine = e
}

func (x *XORM) Add(e *Event) error {
	e.EventId = *new(uint64)
	_, err := x.engine.Table("icode_event").Insert(e)
	if nil != err {
		x.engine.Logger().Errorf("Add event: %v fail : %v", e, err)
		return err
	}
	return nil
}

func (x *XORM) Fetch(eventId uint64) *Event {
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

func (x *XORM) FetchGte(eventId uint64, limit int32) []*Event {
	if int32(0) == limit {
		limit = 20
	}
	return nil
}
