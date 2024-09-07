package event

import (
	"sync"
)

type MemEvent struct {
	maxId  uint64
	events map[uint64]*Event
	mux    *sync.RWMutex
}

func NewMemEvent() *MemEvent {
	m := &MemEvent{}
	m.setup()
	return m
}

func (m *MemEvent) setup() {
	m.maxId = 0
	m.events = make(map[uint64]*Event)
	m.mux = new(sync.RWMutex)
}

func (m *MemEvent) MaxId() uint64 {
	return m.maxId
}

func (m *MemEvent) NextId() uint64 {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.maxId += 1
	return m.maxId
}

func (m *MemEvent) Driver() string {
	return "memory"
}

func (m *MemEvent) Add(e *Event) error {
	//nx := m.NextId()
	m.mux.Lock()
	defer m.mux.Unlock()
	m.maxId += 1
	e.EventId = m.maxId
	m.events[e.EventId] = e
	return nil
}

func (m *MemEvent) Fetch(eventId uint64) *Event {
	m.mux.RLock()
	defer m.mux.RUnlock()
	e, ok := m.events[eventId]
	if ok {
		return e
	}
	return nil
}

func (m *MemEvent) FetchGte(eventId uint64, limit int32) []*Event {
	if int32(0) == limit {
		limit = 20
	}
	return nil
}
