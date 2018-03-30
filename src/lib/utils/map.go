package utils

import (
	"sync"
)

func NewMap() *Map {
	return &Map{elems: map[interface{}]interface{}{}}
}

type Map struct {
	sync.RWMutex
	elems map[interface{}]interface{}
}

func (this *Map) Get(key interface{}) interface{} {
	this.RLock()
	defer this.RUnlock()
	if value, ok := this.elems[key]; ok {
		return value
	} else {
		return nil
	}
}

func (this *Map) Set(key interface{}, value interface{}) bool {
	this.Lock()
	defer this.Unlock()
	if _, ok := this.elems[key]; ok {
		this.elems[key] = value
		return true
	} else {
		this.elems[key] = value
		return false
	}
}

func (this *Map) Del(key interface{}) {
	this.Lock()
	defer this.Unlock()
	delete(this.elems, key)
}

func (this *Map) Len() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.elems)
}

func (this *Map) Exist(key interface{}) bool {
	this.RLock()
	defer this.RUnlock()
	if _, ok := this.elems[key]; ok {
		return true
	} else {
		return false
	}
}

func (this *Map) Range(f func(interface{}, interface{}) bool) {
	this.RLock()
	defer this.RUnlock()
	for k, v := range this.elems {
		if f(k, v) {
			break
		}
	}
}

func (this *Map) LRange(f func(interface{}, interface{}) bool) {
	this.Lock()
	defer this.Unlock()
	for k, v := range this.elems {
		if f(k, v) {
			break
		}
	}
}

func (this *Map) Iter() <-chan interface{} {
	var ch = make(chan interface{})
	go func(s *Map) {
		this.Lock()

		for i := 0; i < len(s.elems); i++ {
			ch <- s.elems[i]
		}

		close(ch)

		this.Unlock()

	}(this)
	return ch
}
