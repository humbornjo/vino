package vino

import (
	"sync"
)

type MutexRW struct {
	lifted bool
	muGate sync.RWMutex
	muYard sync.RWMutex
}

func (m *MutexRW) Lock() {
	m.muGate.Lock()
	m.muYard.Lock()
}

func (m *MutexRW) Unlock() {
	if m.lifted {
		m.lifted = false
		defer m.muGate.RUnlock()
	} else {
		defer m.muGate.Unlock()
	}
	m.muYard.Unlock()
}

func (m *MutexRW) RLock() {
	m.muGate.RLock()
	m.muYard.RLock()
}

func (m *MutexRW) RUnlock() {
	m.muGate.RUnlock()
	m.muYard.RUnlock()
}

func (m *MutexRW) RLifted() {
	m.muGate.RUnlock()
	m.muYard.Lock()
}
