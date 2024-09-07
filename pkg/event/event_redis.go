package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const (
	ICodeEventAutoIncKey = "ic_event_inc_key"
	ICodeEventPrefix     = "ic_event_"
)

type RedisEvent struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisEvent() *RedisEvent {
	return &RedisEvent{}
}

func NewRedisEventWithClient(client *redis.Client) *RedisEvent {
	redis := NewRedisEvent()
	redis.client = client
	redis.ctx = context.Background()
	return redis
}

func (r *RedisEvent) Driver() string {
	return "redis"
}

func (r *RedisEvent) Add(e *Event) error {
	e.EventId, _ = r.client.Incr(r.ctx, ICodeEventAutoIncKey).Uint64()
	key := fmt.Sprintf("%s%v", ICodeEventPrefix, e.EventId)
	data, err := json.Marshal(e)
	if nil != err {
		return err
	}
	r.client.Set(r.ctx, key, data, 0)
	Pub(e)
	return nil
}

func (r *RedisEvent) Fetch(eventId uint64) *Event {
	key := fmt.Sprintf("%s%v", ICodeEventPrefix, eventId)
	e := &Event{}
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		// fmt.Println("err0r", err.Error())
		return nil
	}
	json.Unmarshal(data, e)
	// fmt.Println(string(bytes))
	return e
}
func (r *RedisEvent) FetchGte(eventId uint64, limit int32) []*Event {
	return nil
}
