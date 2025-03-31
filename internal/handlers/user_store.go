package handlers

import (
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	mu    sync.Mutex
	users map[string]string
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]string),
	}
}

func (s *UserStore) AddUser(username, passwordHash string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, err := bcrypt.GenerateFromPassword([]byte(passwordHash), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("hash error: %v\n", err)
		return
	}

	s.users[username] = string(hash)
}

func (s *UserStore) GetPassword(username, password string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, exists := s.users[username]
	if !exists {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserStore) UserExists(username string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.users[username]
	return exists
}
