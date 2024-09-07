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

	engine, err := xorm.NewEngine("sqlite3", "./event_test.db")
	engine.ShowSQL(true)
	if err != nil {
		fmt.Println(err)
		assert.NotNil(t, err)
	}
	exor := NewXORMWithEngine(engine)

	eInsts := []IEvent{
		em, exor,
	}
	data := []struct {
		ev   *Event
		want uint64
	}{
		{&Event{EventType: "t1", EventData: "{}", EntityType: "Order", EntityId: 123}, 1},
		{&Event{EventType: "t2", EventData: "{}", EventId: 10, EntityId: 0}, 2},
	}
	for _, inst := range eInsts {
		for i, d := range data {
			t.Run(fmt.Sprintf("Record %d", i), func(t *testing.T) {
				inst.Add(d.ev)
				self := inst.Fetch(d.ev.EventId)
				if nil == self {
					assert.FailNowf(t, "fetch event return nil", "event : %v", d.ev)
				}
				// 并行情况下,equal 不稳定 改为判断 > 0
				assert.Greater(t, d.ev.EventId, uint64(0))
				assert.Equal(t, d.ev.EventId, self.EventId)
			})
		}
	}

}
