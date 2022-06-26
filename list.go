package procrastiproxy

import (
	"sync"
)

var list = NewList()

type List struct {
	m       sync.Mutex
	members map[string]bool
}

// GetList returns the list singleton
func GetList() *List {
	return list
}

func NewList() *List {
	return &List{
		members: make(map[string]bool),
	}
}

// Clear resets the list, deleting all members
func (l *List) Clear() {
	defer l.m.Unlock()
	l.m.Lock()
	l.members = make(map[string]bool)
}

// All returns every member of the list
func (l *List) All() []string {
	l.m.Lock()
	defer l.m.Unlock()
	var members []string
	for k := range l.members {
		members = append(members, k)
	}
	return members
}

// Add appends an item to the list
func (l *List) Add(item string) {
	l.m.Lock()
	defer l.m.Unlock()
	l.members[item] = true
}

// Remove deletes an item from the list
func (l *List) Remove(item string) {
	l.m.Lock()
	defer l.m.Unlock()
	delete(l.members, item)
}

// Contains returns true if the supplied item is a member of the list
func (l *List) Contains(item string) bool {
	l.m.Lock()
	defer l.m.Unlock()
	return l.members[item]
}

// Length returns the number of members in the list
func (l *List) Length() int {
	l.m.Lock()
	defer l.m.Unlock()
	return len(l.members)
}
