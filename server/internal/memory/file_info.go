package memory

import (
	"fmt"
	"sync"
)

type Hash [32]byte

func NewHash(s string) (h Hash){
	copy(h[:], s)
	return
}

func (h Hash) String() string {
	return string(h[:])
}


type FileInfo struct {
	name string
	h 	Hash
	*sync.Mutex
	holders map[string]User
}

func NewFileInfo(name string, h Hash, holders map[string]User) *FileInfo {
	return &FileInfo{name: name, h: h, holders: holders}
}

func (fi *FileInfo) HasHolder(name string) bool {
	fi.Lock()
	defer fi.Unlock()
	_, exists := fi.holders[name]
	return exists
}

func (fi *FileInfo) HasAnyHolders() bool{
	fi.Lock()
	defer fi.Unlock()
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

func (fi *FileInfo) GetHolders() string {
	fi.Lock()
	defer fi.Unlock()
	return fmt.Sprintf("%v", fi.holders)
}

func (fi *FileInfo) HasSameHash(hash Hash) bool {
	fi.Lock()
	defer fi.Unlock()
	return fi.h == hash
}

