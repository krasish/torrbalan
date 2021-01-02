package memory

import (
	"fmt"
	"sync"
)

type User struct {
	Name string
	Addr string
}

type UserManager struct{
	users map[string]User
	m *sync.Mutex
}

func NewEmptyUserManager() *UserManager {
	return &UserManager{
		users: make(map[string]User),
		m:     &sync.Mutex{},
	}
}

func (um *UserManager) RegisterUser(name, addr string) (User, error) {
	um.m.Lock()
	defer um.m.Unlock()
	if _, exists := um.users[name]; exists {
		return User{}, fmt.Errorf("user %q already exists", name)
	}
	um.users[name] = User{Name: name, Addr: addr}
	return um.users[name], nil
}

func (um *UserManager) DeleteUser(name string) error {
	um.m.Lock()
	defer um.m.Unlock()
	if _, exists := um.users[name]; !exists {
		return fmt.Errorf("user %q does not exists", name)
	}
	delete(um.users, name)
	return nil
}