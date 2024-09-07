package event

import (
	"sync"
)

type Mem struct {
	maxId  uint64
	events map[uint64]*Event
	mux    *sync.RWMutex
}

func NewMem() *Mem {
	m := &Mem{}
	m.setup()
	return m
}

func (m *Mem) setup() {
	m.maxId = 0
	m.events = make(map[uint64]*Event)
	m.mux = new(sync.RWMutex)
}

func (m *Mem) MaxId() uint64 {
	return m.maxId
}

func (m *Mem) NextId() uint64 {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.maxId += 1
	return m.maxId
}

func (m *Mem) Add(e *Event) error {
	//nx := m.NextId()
	m.mux.Lock()
	defer m.mux.Unlock()
	m.maxId += 1
	e.EventId = m.maxId
	m.events[e.EventId] = e
	return nil
}

func (m *Mem) Fetch(pk uint64) *Event {
	m.mux.RLock()
	defer m.mux.RUnlock()
	e, ok := m.events[pk]
	if ok {
		return e
	}
	return nil
}

func (m *Mem) FetchGte(pk uint64, limit int32) []*Event {
	if int32(0) == limit {
		limit = 20
	}
	return nil
}
