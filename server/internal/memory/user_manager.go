package memory

import (
	"fmt"
	"sync"
)

//User represents a single user(client) of the server
//with its username and tcp address.
type User struct {
	Name string
	Addr string
}

//UserManager keeps all the users(clients) registered in the server
//and provides methods for creating and deleting users
//which are safe fo concurrent use.
type UserManager struct {
	users map[string]User
	m     *sync.Mutex
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
