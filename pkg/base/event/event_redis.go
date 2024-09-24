package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/everpan/mdmg/pkg/base/log"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var logger = log.GetLogger()

// icEv = icode event
const (
	ICodeEventAutoIncKey = "icEv_inc_key"
	ICodeEventPrefix     = "icEv_"
)

type RedisEvent struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisEvent() *RedisEvent {
	return &RedisEvent{}
}

func NewRedisEventWithClient(client *redis.Client) *RedisEvent {
	rInst := NewRedisEvent()
	rInst.client = client
	rInst.ctx = context.Background()
	return rInst
}

func (r *RedisEvent) Driver() string {
	return "redis"
}
func (r *RedisEvent) MaxId() uint64 {
	id, err := r.client.Get(r.ctx, ICodeEventAutoIncKey).Uint64()
	if err != nil {
		return 0
	}
	return id
}
func (r *RedisEvent) Add(e *Event) error {
	e.EventId, _ = r.client.Incr(r.ctx, ICodeEventAutoIncKey).Uint64()
	key := fmt.Sprintf("%s%v", ICodeEventPrefix, e.EventId)
	data, err := json.Marshal(e)
	if nil != err {
		return err
	}
	r.client.Set(r.ctx, key, data, 0)
	pubAfter(e)
	return nil
}

func (r *RedisEvent) Fetch(eventId uint64) *Event {
	key := fmt.Sprintf("%s%v", ICodeEventPrefix, eventId)
	e := &Event{}
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		// fmt.Println("err0r", err.Error())
		logger.Error("event get", zap.String("key", key), zap.Error(err))
		return nil
	}
	_ = json.Unmarshal(data, e)
	// fmt.Println(string(bytes))
	return e
}
func (r *RedisEvent) FetchGte(eventId uint64, limit int32) []*Event {
	return nil
}

func (r *RedisEvent) Close() {
	_ = r.client.Close()
}
