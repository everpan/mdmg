package event

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
	"xorm.io/xorm"
)

func TestAddAndFetchEvents(t *testing.T) {
	em := NewMem()

	engine, err := xorm.NewEngine("sqlite3", "./test_event.db")
	engine.ShowSQL(true)
	if err != nil {
		fmt.Println(err)
		assert.NotNil(t, err)
	}
	exor := NewXORMWithEngine(engine)

	eInsts := []IEvent{
		em, exor,
	}
	evs := []struct {
		ev   *Event
		want uint64
	}{
		{&Event{EventType: "t1", EventData: "{}", EntityType: "Order", EntityId: 123}, 1},
		{&Event{EventType: "t2", EventData: "{}", EventId: 10, EntityId: 0}, 2},
	}
	for _, inst := range eInsts {
		for i, e := range evs {
			t.Run(fmt.Sprintf("Record %d", i), func(t *testing.T) {
				inst.Add(e.ev)
				self := inst.Fetch(e.ev.EventId)
				assert.NotNil(t, self)
				assert.Equal(t, e.ev.EventId, self.EventId)
				assert.Greater(t, e.ev.EventId, uint64(0))
			})
		}
	}

}
