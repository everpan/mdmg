package event

import (
	"sync"
)

type EventMem struct {
	maxId  uint64
	events map[uint64]*Event
	mux    *sync.RWMutex
}

func (em *EventMem) Init() {
	em.maxId = 0
	em.events = make(map[uint64]*Event)
	em.mux = new(sync.RWMutex)
}

func (em *EventMem) MaxId() uint64 {
	return em.maxId
}

func (em *EventMem) NextId() uint64 {
	em.mux.Lock()
	defer em.mux.Unlock()
	em.maxId += 1
	return em.maxId
}

func (em *EventMem) Add(e *Event) error {
	//nx := em.NextId()
	em.mux.Lock()
	defer em.mux.Unlock()
	em.maxId += 1
	e.EventID = em.maxId
	em.events[e.EventID] = e
	return nil
}

func (em *EventMem) Fetch(pk uint64) *Event {
	em.mux.RLock()
	defer em.mux.RUnlock()
	e, ok := em.events[pk]
	if ok {
		return e
	}
	return nil
}

func (ex *EventMem) FetchGte(pk uint64, limit int32) []*Event {
	if int32(0) == limit {
		limit = 20
	}
	return nil
}
