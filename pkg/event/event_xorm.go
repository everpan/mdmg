package event

import (
	"xorm.io/xorm"
)

type EventXORM struct {
	engine *xorm.Engine
}

func (ex *EventXORM) InitTables() {
	ex.engine.Table("icode_event").CreateTable(&Event{})
	ex.engine.Table("icode_entity").CreateTable(&Entity{})
	ex.engine.Table("icode_snapshot").CreateTable(&Snapshot{})
}

func (ex *EventXORM) SetEngine(e *xorm.Engine) {
	ex.engine = e
}

func (ex *EventXORM) Add(e *Event) error {
	e.EventID = *new(uint64)
	_, err := ex.engine.Insert(e)
	if nil != err {
		ex.engine.Logger().Errorf("Add event: %v fail : %v", e, err)
		return err
	}
	return nil
}

func (ex *EventXORM) Fetch(pk uint64) *Event {
	return nil
}

func (ex *EventXORM) FetchGte(pk uint64, limit int32) []*Event {
	if int32(0) == limit {
		limit = 20
	}
	return nil
}
