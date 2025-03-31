package vino

import (
	"sync"
)

// MutexRW is a RWMutex that can be upgrade (a.k.a lifted)
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

// "Rpgrade" pronounces just like "Upgrade", so coooool
func (m *MutexRW) Rpgrade() {
	m.muYard.RUnlock()
	m.muYard.Lock()
	m.lifted = true
}

func (m *MutexRW) Degrade() {
	m.muYard.Unlock()
	m.muYard.RLock()
	m.lifted = false
}

func (m *MutexRW) TryLock() bool {
	if m.muGate.TryLock() {
		if m.muYard.TryLock() {
			return true
		}
		m.muGate.Unlock()
	}
	return false
}

func (m *MutexRW) TryRLock() bool {
	if m.muGate.TryRLock() {
		if m.muYard.TryRLock() {
			return true
		}
		m.muGate.RUnlock()
	}
	return false
}
