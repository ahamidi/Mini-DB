package main

import (
	"errors"
	"strconv"
	"sync"
)

type MemDB struct {
	Records  map[string]*Record
	Mutex    sync.Mutex // DB-wide lock
	Autoincr int        // Using counter as lock_id (assuming security is not an issue)
}

type Record struct {
	Value interface{} `json:"value,omitempty"` // omitempty in the case of overriding the response for create key handler
	Lock  string      `json:"lock_id"`
	Mutex sync.Mutex  `json:"-"` // per record lock
}

func newMemDB() *MemDB {
	recs := make(map[string]*Record)
	var m sync.Mutex

	return &MemDB{recs, m, 1}
}

func (m *MemDB) NextLockID() string {
	m.Autoincr++
	next := strconv.Itoa(m.Autoincr)
	return next
}

// Retrieve value of {key}. Also locks key and returns `lock_id`
func (m *MemDB) getRecordAndLock(key string) (*Record, error) {
	if rec, ok := m.Records[key]; ok {
		rec.Mutex.Lock()
		rec.Lock = m.NextLockID()
		return rec, nil
	}
	return nil, errors.New("Key not found")
}

// Sets the value of a key provided lock ID matches. Also supports releasing lock.
func (m *MemDB) updateKey(key string, lockID string, value interface{}, release bool) error {

	if rec, ok := m.Records[key]; !ok {
		return errors.New("Key Not Found")
	} else {
		// Check if lock ID matches
		if lockID == rec.Lock {
			rec.Value = value
		} else {
			return errors.New("Unauthorized")
		}

		if release {
			rec.Mutex.Unlock()
			rec.Lock = ""
		}
		return nil
	}

	return nil
}

// If key exists, it updates the value stored, otherwise creats new record.
func (m *MemDB) createKey(key string, value interface{}) (*Record, error) {
	// Check if key already exists
	rec, err := m.getRecordAndLock(key)
	if err != nil { // Rec doesn't exist, so create it
		m.Mutex.Lock()
		defer m.Mutex.Unlock()
		rec = &Record{
			Value: value,
			Lock:  m.NextLockID(),
		}
		rec.Mutex.Lock()
		m.Records[key] = rec
	}
	return rec, nil
}
