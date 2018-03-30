package utils

import "sync"

type Array interface {
	//Push 向链表尾部添加一个或者多个元素
	Push(values ...interface{})

	// Pop 从列表尾部移除并返回一个元素
	Pop() interface{}

	//PushFront 向链表头部添加一个元素
	PushFront(values ...interface{})

	// PopFront从链表头部移除并返回一个元素
	PopFront() interface{}

	// Len 返回集链表元素的长度
	Len() int

	// Values 返回链表所有元素组成的 Slice
	Values() []interface{}

	// Iter 遍历链表的所有元素
	Iter() <-chan interface{}

	// Clear  清空链表
	Clear()
}

type array struct {
	m     []interface{}
	rw    sync.RWMutex
	block bool
}

// NewArray 创建一个链表，block == true :线程安全，block == false:非线程安全
func NewArray(block bool, values ...interface{}) Array {
	var s = &array{}
	s.block = block
	if len(values) > 0 {
		s.Push(values...)
	}
	return s
}

func (this *array) lock() {
	if this.block {
		this.rw.Lock()
	}
}

func (this *array) unlock() {
	if this.block {
		this.rw.Unlock()
	}
}

func (this *array) rLock() {
	if this.block {
		this.rw.RLock()
	}
}

func (this *array) rUnlock() {
	if this.block {
		this.rw.RUnlock()
	}
}

func (this *array) Push(values ...interface{}) {
	this.lock()
	defer this.unlock()
	this.m = append(this.m, values...)
}

func (this *array) PushFront(values ...interface{}) {
	this.lock()
	defer this.unlock()
	var tmp []interface{}
	tmp = append(tmp, values...)
	this.m = append(tmp, this.m...)
}

func (this *array) PopFront() interface{} {
	this.lock()
	defer this.unlock()
	if this.len() > 0 {
		elem := this.m[0]
		this.m = this.m[1:]
		return elem
	}
	return nil
}

func (this *array) Pop() interface{} {
	this.lock()
	defer this.unlock()
	if this.len() > 0 {
		elem := this.m[this.len()-1]
		this.m = this.m[:this.len()-1]
		return elem
	}
	return nil
}

func (this *array) Len() int {
	this.rLock()
	defer this.rUnlock()

	return this.len()
}

func (this *array) len() int {
	return len(this.m)
}

func (this *array) Values() []interface{} {
	this.rLock()
	defer this.rUnlock()
	tmp := make([]interface{}, this.len())
	if this.len() > 0 {
		copy(tmp, this.m)
	}
	return tmp
}

func (this *array) Iter() <-chan interface{} {
	var ch = make(chan interface{})

	go func(s *array) {
		if s.block {
			s.rLock()
		}

		for i := 0; i < s.len(); i++ {
			ch <- this.m[i]
		}

		close(ch)

		if s.block {
			s.rUnlock()
		}

	}(this)

	return ch
}

func (this *array) Clear() {
	this.m = []interface{}{}
}
