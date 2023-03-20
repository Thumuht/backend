package db

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/uptrace/bun"
)

// Simple Test for Mapcache
func TestMapCacheSetGet(t *testing.T) {
	c := NewMapCache[string, int]()
	c.Set("hello", 123)
	c.Set("world, 456", 456)
	v, ok := c.Get("hello")
	if v == nil || *v != 123 || ok != true {
		t.Errorf("set cache wrong")
	}

	v, ok = c.Get("world, 456")
	if v == nil || *v != 456 || ok != true {
		t.Errorf("set cache wrong")
	}

	c.Remove("world, 456")
	v, ok = c.Get("world, 456")
	if ok == true || v != nil {
		t.Errorf("delete cache wrong")
	}

	c.Invalidate(false)
	v, ok = c.Get("hello")
	if ok == true || v != nil {
		t.Errorf("invalidate cache wrong")
	}
}

// Concurrency Test for Mapcache
func TestMapCacheConcurrency(t *testing.T) {
	c := MapCache[string, int]{
		c: make(map[string]int, 0),
	}
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		// i must be in parameter
		go func(i int) {
			defer wg.Done()
			c.Set("Test", i)
		}(i)
	}

	wg.Wait()
	v, ok := c.Get("Test")
	if !ok {
		t.Errorf("concurrency problem %v %v", ok, *v)
	}
}

// Test Flush
func TestMapCacheFlush(t *testing.T) {
	ctx := context.Background()
	testpdb, err := InitSQLiteDB()
	if err != nil {
		t.Error("cannot new db")
	}
	mc := NewMapCache[string, int]()
	mc.SetFlushTarget("test_like", "name", "like", testpdb)

	// test relation
	type TestLike struct {
		bun.BaseModel `bun:"table:test_like"`

		Name string `bun:"name"`
		Like int32  `bun:"like"`
	}

	_, err = testpdb.NewCreateTable().Model((*TestLike)(nil)).Exec(ctx)
	if err != nil {
		t.Error("cannot create table")
	}

	mc.Set("hello", 123)
	mc.Set("world, 456", 456)

	mc.Invalidate(true)

	var dbcache []TestLike
	err = testpdb.NewSelect().Model(&dbcache).Scan(ctx)
	if err != nil {
		t.Error("cannot query table")
	}
	fmt.Println(dbcache)
}
