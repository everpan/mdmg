package event

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndFetchEvents(t *testing.T) {
	em := new(EventMem)
	em.Init()
	eInsts := []EventInterface{
		em,
	}
	evs := []struct {
		ev   *Event
		want uint64
	}{
		{&Event{EventType: "t1", EventData: "{}", EntityType: "Order", EntityId: 123}, 1},
		{&Event{EventType: "t2", EventData: "{}", EventID: 10, EntityId: 0}, 2},
	}
	for _, inst := range eInsts {
		for i, e := range evs {
			t.Run(fmt.Sprintf("Record %d", i), func(t *testing.T) {
				inst.Add(e.ev)
				self := inst.Fetch(e.ev.EventID)
				assert.Equal(t, e.ev, self)
				assert.Equal(t, e.want, e.ev.EventID)
			})
		}
	}

}
