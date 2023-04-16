package db

import (
	"fmt"
	"strings"

	"github.com/uptrace/bun"
)

// This file implements a memcache using CSP (Communicating Sequential Processes)

// MapCacheCSP is a concurrent-safe map cache
type MapCacheCSP[K comparable, V any] struct {
	c      map[K]V
	table  string
	keycol string
	valcol string
	pdb    *bun.DB // persistant storage db
	// channels
	eventc chan event
	reqc   chan []any
	quitc  chan struct{}
}

// event represents Req/Resp
type enumEventType int

const (
	GET enumEventType = iota
	REMOVE
	SET_VALUE
	SET_TARGET
	FLUSH
	STRING
	INVALIDATE
)

type event struct {
	eventType enumEventType
	args      []any
}

func (mc *MapCacheCSP[K, V]) loop() {
	for {
		select {
		case e := <-mc.eventc:
			mc.reqc <- mc.process(e)
		case <-mc.quitc:
			close(mc.reqc)
			close(mc.eventc)
			close(mc.quitc)
			return
		}
	}
}

// string is a single-threaded String method
func (mc *MapCacheCSP[K, V]) string() string {
	var s []string
	for k, v := range mc.c {
		s = append(s, fmt.Sprintf("('%#v', '%#v')", k, v))
	}
	return strings.Join(s, ",\n")
}

// single-threaded flush()
func (mc *MapCacheCSP[K, V]) flush() error {
	stat := fmt.Sprintf(`
	REPLACE INTO %s (%s, %s)
	VALUES %s;
	`, mc.table, mc.keycol, mc.valcol, mc.string())

	_, err := mc.pdb.Exec(stat)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MapCacheCSP[K, V]) process(e event) []any {
	switch e.eventType {
	case GET:
		value, ok := mc.c[e.args[0].(K)]
		if !ok {
			return []any{
				nil,
				false,
			}
		}
		return []any{
			&value,
			ok,
		}
	case REMOVE:
		delete(mc.c, e.args[0].(K))
	case SET_VALUE:
		mc.c[e.args[0].(K)] = e.args[1].(V)
	case SET_TARGET:
		mc.table = e.args[0].(string)
		mc.keycol = e.args[1].(string)
		mc.valcol = e.args[2].(string)
		mc.pdb = e.args[3].(*bun.DB)
	case FLUSH:
		err := mc.flush()
		if err != nil {
			return []any{
				err,
			}
		}
	case STRING:
		var s []string
		for k, v := range mc.c {
			s = append(s, fmt.Sprintf("('%#v', '%#v')", k, v))
		}
		return []any{
			strings.Join(s, ",\n"),
		}
	case INVALIDATE:
		if e.args[0].(bool) {
			err := mc.flush()
			if err != nil {
				return []any{
					err,
				}
			}
		}
		for k := range mc.c {
			delete(mc.c, k)
		}
	}
	return nil
}

func (mc *MapCacheCSP[K, V]) Get(key K) (*V, bool) {
	mc.eventc <- event{
		GET,
		[]any{
			key,
		},
	}
	resp := <-mc.reqc
	if resp[0] == nil {
		return nil, resp[1].(bool)
	}
	return resp[0].(*V), resp[1].(bool)
}

func (mc *MapCacheCSP[K, V]) Remove(key K) error {
	mc.eventc <- event{
		REMOVE,
		[]any{
			key,
		},
	}
	<-mc.reqc
	return nil
}

func (mc *MapCacheCSP[K, V]) Set(key K, value V) error {
	mc.eventc <- event{
		SET_VALUE,
		[]any{
			key,
			value,
		},
	}
	<-mc.reqc
	return nil
}

func (mc *MapCacheCSP[K, V]) SetFlushTarget(table, keycol, valcol string, pdb *bun.DB) {
	mc.eventc <- event{
		SET_TARGET,
		[]any{
			table,
			keycol,
			valcol,
			pdb,
		},
	}
	<-mc.reqc
}

func (mc *MapCacheCSP[K, V]) Flush() error {
	mc.eventc <- event{
		FLUSH,
		nil,
	}
	resp := <-mc.reqc
	if resp != nil {
		return resp[0].(error)
	}
	return nil
}

func (mc *MapCacheCSP[K, V]) String() string {
	mc.eventc <- event{
		STRING,
		nil,
	}
	resp := <-mc.reqc
	return resp[0].(string)
}

func (mc *MapCacheCSP[K, V]) Invalidate(flush bool) error {
	mc.eventc <- event{
		INVALIDATE,
		[]any{
			flush,
		},
	}
	resp := <-mc.reqc
	if resp != nil {
		return resp[0].(error)
	}
	return nil
}

func (mc *MapCacheCSP[K, V]) Shutdown() {
	mc.quitc <- struct{}{}
}

func NewMapCacheCSP[K comparable, V any]() Cache[K, V] {
	newmc := &MapCacheCSP[K, V]{
		c:      make(map[K]V),
		eventc: make(chan event),
		reqc:   make(chan []any),
		quitc:  make(chan struct{}),
	}
	go newmc.loop()
	return newmc
}
