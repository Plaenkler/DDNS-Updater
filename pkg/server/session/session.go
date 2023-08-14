package session

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

const (
	sessionDuration = time.Minute * 5
	houseKeepingInt = time.Minute * 1
)

var (
	sessions = make(map[string]time.Time)
	oc       sync.Once
	mu       sync.Mutex
	cancel   = make(chan bool)
)

func Start() {
	oc.Do(func() {
		go startCaretaker()
	})
}

func startCaretaker() {
	ticker := time.NewTicker(houseKeepingInt)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cleanOrphans()
			continue
		case <-cancel:
			return
		}
	}
}

func cleanOrphans() {
	currentTime := time.Now()
	mu.Lock()
	defer mu.Unlock()
	for key, s := range sessions {
		if s.Before(currentTime) {
			delete(sessions, key)
		}
	}
}

func Add() (string, error) {
	key, err := generateKey()
	if err != nil {
		return "", err
	}
	mu.Lock()
	defer mu.Unlock()
	sessions[key] = time.Now().Add(sessionDuration)
	return key, nil
}

func generateKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(key), nil
}

func Verify(key string) bool {
	mu.Lock()
	defer mu.Unlock()
	s, ok := sessions[key]
	if !ok {
		return false
	}
	if s.Before(time.Now()) {
		delete(sessions, key)
		return false
	}
	sessions[key] = time.Now().Add(sessionDuration)
	return true
}

func Stop() {
	cancel <- true
}
