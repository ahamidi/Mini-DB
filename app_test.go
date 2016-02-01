package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateKeyHandler(t *testing.T) {
	// Init MemDB
	db := newMemDB()
	context.db = db

	req, _ := http.NewRequest("PUT", "/values/a", strings.NewReader("hello"))
	w := httptest.NewRecorder()
	handlers().ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"lock_id":"2"}`, w.Body.String())
}

func TestResNotFound(t *testing.T) {
	// Init MemDB
	db := newMemDB()
	context.db = db

	req, _ := http.NewRequest("POST", "/reservations/a", nil)
	w := httptest.NewRecorder()
	handlers().ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestRes(t *testing.T) {
	// Init MemDB
	db := newMemDB()
	context.db = db

	// Seed DB
	rec := &Record{
		Value: "hello",
		Lock:  "100",
	}
	db.Records["a"] = rec

	req, _ := http.NewRequest("POST", "/reservations/a", nil)
	w := httptest.NewRecorder()
	handlers().ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"value":"hello","lock_id":"2"}`, w.Body.String())

	// Make sure record mutex is locked
	m := sync.Mutex{}
	m.Lock()
	assert.Equal(t, m, db.Records["a"].Mutex)
}

func TestUpdateNotFound(t *testing.T) {
	// Init MemDB
	db := newMemDB()
	context.db = db

	// Seed DB
	rec := &Record{
		Value: "hello",
		Lock:  "555",
	}
	db.Records["a"] = rec

	req, _ := http.NewRequest("POST", "/values/b/555", strings.NewReader("goodbye"))
	w := httptest.NewRecorder()
	handlers().ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestUpdateNotAuthorized(t *testing.T) {
	// Init MemDB
	db := newMemDB()
	context.db = db

	// Seed DB
	rec := &Record{
		Value: "hello",
		Lock:  "555",
	}
	db.Records["a"] = rec

	req, _ := http.NewRequest("POST", "/values/a/444", strings.NewReader("goodbye"))
	w := httptest.NewRecorder()
	handlers().ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestUpdateRelease(t *testing.T) {
	// Init MemDB
	db := newMemDB()
	context.db = db

	// Seed DB
	rec := &Record{
		Value: "hello",
		Lock:  "555",
	}
	rec.Mutex.Lock()
	db.Records["a"] = rec

	req, _ := http.NewRequest("POST", "/values/a/555?release=true", strings.NewReader("goodbye"))
	w := httptest.NewRecorder()
	handlers().ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)

	// Test mutex is unlocked
	m := sync.Mutex{}
	assert.Equal(t, m, db.Records["a"].Mutex)
}

func TestUpdateNoRelease(t *testing.T) {
	// Init MemDB
	db := newMemDB()
	context.db = db

	// Seed DB
	rec := &Record{
		Value: "hello",
		Lock:  "555",
	}
	rec.Mutex.Lock()
	db.Records["a"] = rec

	req, _ := http.NewRequest("POST", "/values/a/555?release=false", strings.NewReader("goodbye"))
	w := httptest.NewRecorder()
	handlers().ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)

	// Test mutex is unlocked
	m := sync.Mutex{}
	m.Lock()
	assert.Equal(t, m, db.Records["a"].Mutex)
}
