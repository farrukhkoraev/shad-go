//go:build !solution

package lrucache

import "container/list"
// import "fmt"

type pair struct {
	first, second int
}

type LruCache struct {
	capacity int
	keyMap map[int](*list.Element)
	accessList *list.List
}

func (c LruCache) Get(key int) (int, bool) {
	e, ok := c.keyMap[key]
	if ok {
		c.accessList.MoveToFront(e)
		p := e.Value.(pair)
		return p.second, true
		
	} else {
		return 0, false
	}
}

func (c LruCache) Set(key, value int) {
	if c.capacity == 0 {
		return 
	}
	
	e, ok := c.keyMap[key]
	
	if ok {
		e.Value = pair{key, value}
		c.accessList.MoveToFront(e)
	} else {
		
		// full capacity, remove lru  `Element` from `accessList` and `acessMap`
		
		if c.accessList.Len() == c.capacity {
			last := c.accessList.Back()
		    kv := last.Value.(pair)
			
		    delete(c.keyMap, kv.first)
			c.accessList.Remove(last)
		}
		
		e = c.accessList.PushFront(pair{key, value})
		c.keyMap[key] = e
	}

}

func (c LruCache) Range(f func(key, value int) bool) {
	
	for e := c.accessList.Back(); e != nil; e = e.Prev() {
		kv := e.Value.(pair)
		if !f(kv.first, kv.second) {
			break
		}
	}
}

func (c LruCache) Clear () {
	for k := range c.keyMap {
		delete(c.keyMap, k)
	}
	c.accessList = c.accessList.Init()
}

func New(cap int) Cache {
	return LruCache {
		capacity: cap,
		keyMap: make(map[int](*list.Element), cap),
		accessList: list.New(),
	}
}
