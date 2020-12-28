package memory

import (
	"fmt"
)

type User struct {
	name string
	addr string
}

type UserManager struct{
	users map[string]User
}

func NewEmptyUserManager() *UserManager {
	return &UserManager{users: make(map[string]User)}
}

func (um *UserManager) RegisterUser(name, addr string) (User, error) {
	if _, exists := um.users[name]; exists {
		return User{}, fmt.Errorf("user %q already exists", name)
	}
	um.users[name] = User{name: name, addr: addr}
	return um.users[name], nil
}

func (um *UserManager) DeleteUser(name string) error {
	if _, exists := um.users[name]; !exists {
		return fmt.Errorf("user %q does not exists", name)
	}
	delete(um.users, name)
	return nil
}