package utils

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	var s = NewArray(true, 1, 2, 3, 5, 6, 7, 7)

	t.Error("集合的长度应该为 7", s.Len())

	s.Push(111, 444)
	s.PushFront(222, 555)
	e := s.Pop()
	fmt.Println(e)
	e = s.PopFront()
	fmt.Println(e)
	fmt.Println(s.Values())
}

func TestIter(t *testing.T) {
	fmt.Println("=====TestIter=====")
	var s1 = NewArray(true, 1, 2, 3, 4, 5)

	for v := range s1.Iter() {
		fmt.Println(v)
	}
}
