package db

import (
	"fmt"
	"log"

	"github.com/uptrace/bun"
)

// cache interface hides the actual cache service from user
//
// uses generic magic
type Cache[K comparable, V any] interface {
	Get(key K) (*V, bool)
	Remove(key K) error
	Set(key K, value V) error
	SetFlushTarget(table, keycol, valcol string, pdb bun.DB)
	Flush() error
	Invalidate() error
}

// cache group manages all cache
type CacheGroup map[string]any // any must be a cache!!

// IN MEMORY MAP CACHE
type MapCache[K comparable, V any] struct {
	c      map[K]V
	table  string
	keycol string
	valcol string
	pdb    bun.DB // persistant storage db
}

func (mc *MapCache[K, V]) Get(key K) (*V, bool) {
	value, ok := mc.c[key]
	if !ok {
		return nil, false
	}
	return &value, true
}

func (mc *MapCache[K, V]) Remove(key K) error {
	delete(mc.c, key)
	return nil
}

func (mc *MapCache[K, V]) Set(key K, value V) error {
	mc.c[key] = value
	return nil
}

func (mc *MapCache[K, V]) SetFlushTarget(table, keycol, valcol string, pdb bun.DB) {
	mc.table = table
	mc.keycol = keycol
	mc.valcol = valcol
	mc.pdb = pdb
}

func (mc *MapCache[K, V]) Flush() error {
	stat := fmt.Sprintf(`WITH _data (%s, %s) AS (
		VALUES
%s
	)
	UPDATE %s
	SET %s = _data.%s
	FROM _data
	WHERE %s.%s = _data.%s;
	`, mc.keycol, mc.valcol, mc.String(), mc.table, mc.valcol, mc.valcol, mc.table, mc.keycol, mc.keycol)

	log.Println(stat)
	_, err := mc.pdb.Exec(stat)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MapCache[K, V]) String() string {
	var s string
	for k, v := range mc.c {
		s += fmt.Sprintf("('%#v', '%#v')\n", k, v)
	}
	return s
}

func (mc *MapCache[K, V]) Invalidate() error {
	for k := range mc.c {
		delete(mc.c, k)
	}
	return nil
}

func NewMapCache[K comparable, V any]() Cache[K, V] {
	return &MapCache[K, V]{
		c: make(map[K]V),
	}
}
