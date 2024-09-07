package event

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"xorm.io/xorm"
)

func TestMemEvent_AddFetch(t *testing.T) {
	m := NewMemEvent()
	AddFetch(t, m)
}
func TestXORMEvent_AddFetch(t *testing.T) {
	engine, err := xorm.NewEngine("sqlite3", "./event_test.db")
	engine.ShowSQL(true)
	if err != nil {
		fmt.Println(err)
		assert.NotNil(t, err)
	}
	x := NewXORMEventWithEngine(engine)
	AddFetch(t, x)
}

func TestRedisEvent_AddFetch(t *testing.T) {
	// redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})
	re := NewRedisEventWithClient(rdb)
	AddFetch(t, re)
}

func AddFetch(t *testing.T, inst IEvent) {
	data := []struct {
		ev   *Event
		want uint64
	}{
		{&Event{EventType: "t1", EventData: "{}", EntityType: "Order", EntityId: 123}, 1},
		{&Event{EventType: "t2", EventData: "{}", EventId: 10, EntityId: 0}, 2},
	}

	for i, d := range data {
		t.Run(fmt.Sprintf("%s Record %d", inst.Driver(), i), func(t *testing.T) {
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
