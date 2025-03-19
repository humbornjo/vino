package vino

import (
	"sync"
)

type RWMutex struct {
	lifted bool
	john   sync.RWMutex
	jane   sync.RWMutex
}

func (m *RWMutex) Lock() {
	m.john.Lock()
	m.jane.Lock()
}

func (m *RWMutex) Unlock() {
	if m.lifted {
		m.lifted = false
		defer m.john.RUnlock()
	} else {
		defer m.john.Unlock()
	}
	m.jane.Unlock()
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
}
