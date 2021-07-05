package adapter

import (
	"github.com/robertscherbarth/go-service-skeleton/internal/users"
	"sync"
)

type InMemoryStore struct {
	sync.RWMutex
	store map[string]users.User
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[string]users.User, 0),
	}
}

func (i *InMemoryStore) Create(user users.User) error {
	i.RWMutex.Lock()
	defer i.RWMutex.Unlock()

	i.store[user.ID.String()] = user
	return nil
}

func (i *InMemoryStore) Delete(id string) error {
	i.RWMutex.Lock()
	defer i.RWMutex.Unlock()

	delete(i.store, id)
	return nil
}

func (i *InMemoryStore) FindByID(id string) (users.User, error) {
	i.RWMutex.RLock()
	defer i.RWMutex.RUnlock()

	return i.store[id], nil
}

func (i *InMemoryStore) FindAll() ([]users.User, error) {
	i.RWMutex.RLock()
	defer i.RWMutex.RUnlock()

	var users []users.User
	for _, v := range i.store {
		users = append(users, v)
	}
	return users, nil
}
