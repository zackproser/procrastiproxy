package procrastiproxy

import (
	"fmt"
	"sync"
)

var list *List

type List struct {
	m       sync.Mutex
	members map[string]bool
}

func GetList() *List {
	if list != nil {
		return list
	}
	list = NewList()
	return list
}

func NewList() *List {
	return &List{
		members: make(map[string]bool),
	}
}

func (l *List) Dump() {
	fmt.Println("List dumping...")
	l.m.Lock()
	defer l.m.Unlock()
	for k, v := range l.members {
		fmt.Printf("key: %v - value: %v\n", k, v)
	}
}

func (l *List) All() []string {
	l.m.Lock()
	defer l.m.Unlock()
	var members []string
	for k := range l.members {
		members = append(members, k)
	}
	return members
}

func (l *List) Add(item string) {
	l.m.Lock()
	defer l.m.Unlock()
	l.members[item] = true
}

func (l *List) Remove(item string) {
	l.m.Lock()
	defer l.m.Unlock()
	delete(l.members, item)
}

func (l *List) Contains(item string) bool {
	l.m.Lock()
	defer l.m.Unlock()
	return l.members[item]
}
