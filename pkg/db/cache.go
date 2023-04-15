package db

import (
	"fmt"
	"strings"
	"sync"

	"github.com/uptrace/bun"
)

// cache interface hides the actual cache service from user
//
// uses generic magic
type Cache[K comparable, V any] interface {
	Get(key K) (*V, bool)
	Remove(key K) error
	Set(key K, value V) error
	SetFlushTarget(table, keycol, valcol string, pdb *bun.DB)
	Flush() error
	Invalidate(flush bool) error
	Shutdown()
}

// IN MEMORY MAP CACHE
type MapCache[K comparable, V any] struct {
	c            map[K]V
	table        string
	keycol       string
	valcol       string
	pdb          *bun.DB // persistant storage db
	sync.RWMutex         // Read-Write lock. Must Get it before r/w global state.
}

func (mc *MapCache[K, V]) Get(key K) (*V, bool) {
	mc.RLock()
	defer mc.RUnlock()
	value, ok := mc.c[key]
	if !ok {
		return nil, false
	}
	return &value, true
}

func (mc *MapCache[K, V]) Remove(key K) error {
	mc.Lock()
	defer mc.Unlock()
	delete(mc.c, key)
	return nil
}

func (mc *MapCache[K, V]) Set(key K, value V) error {
	mc.Lock()
	defer mc.Unlock()
	mc.c[key] = value
	return nil
}

func (mc *MapCache[K, V]) SetFlushTarget(table, keycol, valcol string, pdb *bun.DB) {
	mc.Lock()
	defer mc.Unlock()
	mc.table = table
	mc.keycol = keycol
	mc.valcol = valcol
	mc.pdb = pdb
}

func (mc *MapCache[K, V]) Flush() error {
	stat := fmt.Sprintf(`
	REPLACE INTO %s (%s, %s)
	VALUES %s;
	`, mc.table, mc.keycol, mc.valcol, mc.String())

	_, err := mc.pdb.Exec(stat)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MapCache[K, V]) String() string {
	mc.RLock()
	defer mc.RUnlock()
	var s []string
	for k, v := range mc.c {
		s = append(s, fmt.Sprintf("('%#v', '%#v')", k, v))
	}
	return strings.Join(s, ",\n")
}

func (mc *MapCache[K, V]) Invalidate(flush bool) error {
	if flush {
		mc.Flush()
	}
	mc.Lock()
	defer mc.Unlock()
	for k := range mc.c {
		delete(mc.c, k)
	}
	return nil
}

func (mc *MapCache[K, V]) Shutdown() {
	mc.Invalidate(true)
}

func NewMapCache[K comparable, V any]() Cache[K, V] {
	return &MapCache[K, V]{
		c: make(map[K]V),
	}
}
