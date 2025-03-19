package vino

import "sync"

type RWMutex struct {
	upgraded bool
	john     sync.RWMutex
	jane     sync.RWMutex
}

func (m *RWMutex) Lock() {
	m.john.Lock()
	m.jane.Lock()
}

func (m *RWMutex) Unlock() {
	m.jane.Unlock()
	if m.upgraded {
		m.upgraded = false
		m.john.RUnlock()
		return
	}
	m.john.Unlock()
}

func (m *RWMutex) RLock() {
	m.john.RLock()
	m.jane.RLock()
}

func (m *RWMutex) RUnlock() {
	m.jane.RUnlock()
	m.john.RUnlock()
}

func (m *RWMutex) RLifted() {
	m.jane.RUnlock()
	m.jane.Lock()
	m.upgraded = true
}
