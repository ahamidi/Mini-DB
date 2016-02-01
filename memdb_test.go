package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRecordAndLock(t *testing.T) {
	db := newMemDB()

	val := "hello"
	db.Records["a"] = &Record{
		Value: val,
	}

	db.Records["b"] = &Record{
		Value: "hi",
	}
	rec, err := db.getRecordAndLock("a")

	assert.NoError(t, err)
	assert.NotNil(t, rec)
	assert.Equal(t, val, rec.Value)
}

func TestCreateKey(t *testing.T) {
	db := newMemDB()

	rec, err := db.createKey("aaa", "hello")

	assert.NoError(t, err)
	assert.NotNil(t, rec)
	assert.Equal(t, "hello", rec.Value)

	// Check that the mutex is locked by function
	m := sync.Mutex{}
	m.Lock()
	assert.Equal(t, rec.Mutex, m)
}

func TestUpdateKey(t *testing.T) {
	db := newMemDB()

	rec := &Record{
		Value: "first",
		Lock:  "100",
	}
	rec.Mutex.Lock()
	db.Records["bbb"] = rec

	err := db.updateKey("bbb", "100", "second", false)

	assert.NoError(t, err)
	assert.Equal(t, "second", rec.Value)
}

func TestUpdateKeyUnauthorized(t *testing.T) {
	db := newMemDB()

	rec := &Record{
		Value: "first",
		Lock:  "100",
	}
	rec.Mutex.Lock()
	db.Records["bbb"] = rec

	err := db.updateKey("bbb", "1", "second", false)
	assert.Error(t, err)
	assert.Equal(t, "Unauthorized", err.Error())
}

func TestUpdateKeyNotFound(t *testing.T) {
	db := newMemDB()

	err := db.updateKey("ccc", "1", "second", false)
	assert.Error(t, err)
	assert.Equal(t, "Key Not Found", err.Error())
}
