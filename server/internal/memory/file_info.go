package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type Hash [64]byte

func NewHash(s string) (h Hash) {
	copy(h[:], s)
	return
}

func (h Hash) String() string {
	return string(h[:])
}

type FileInfo struct {
	name string
	h    Hash
	*sync.RWMutex
	holders map[string]User
}

func NewFileInfo(name string, h Hash, holders map[string]User) *FileInfo {
	return &FileInfo{name: name, h: h, holders: holders, RWMutex: &sync.RWMutex{}}
}

func (fi *FileInfo) HasHolder(name string) bool {
	fi.RLock()
	defer fi.RUnlock()
	_, exists := fi.holders[name]
	return exists
}

func (fi *FileInfo) HasAnyHolders() bool {
	fi.RLock()
	defer fi.RUnlock()
	return len(fi.holders) == 0
}

func (fi *FileInfo) AddHolder(user User) error {
	if fi.HasHolder(user.Name) {
		return fmt.Errorf("user %q already is holder of file %s", user.Name, fi.name)
	}
	fi.Lock()
	defer fi.Unlock()
	fi.holders[user.Name] = user
	return nil
}

func (fi *FileInfo) RemoveHolder(username string) error {
	if !fi.HasHolder(username) {
		return fmt.Errorf("user %q is not a holder of file %s", username, fi.name)
	}
	fi.Lock()
	defer fi.Unlock()
	delete(fi.holders, username)
	return nil
}

func (fi *FileInfo) GetHolders() ([]byte, error) {
	fi.RLock()
	defer fi.RUnlock()
	holdersSlice, i := make([]User, len(fi.holders)), 0
	for _, user := range fi.holders {
		holdersSlice[i] = user
		i++
	}
	marshaledHolders, err := json.Marshal(holdersSlice)
	if err != nil {
		return nil, errors.New("could not marshal holders")
	}
	return marshaledHolders, nil

}

func (fi *FileInfo) HasSameHash(hash Hash) bool {
	fi.RLock()
	defer fi.RUnlock()
	return fi.h == hash
}
