package db

import (
	"testing"
)

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

	c.Invalidate()
	v, ok = c.Get("hello")
	if ok == true || v != nil {
		t.Errorf("invalidate cache wrong")
	}
}
