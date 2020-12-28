package memory

import (
	"fmt"
	"log"
	"sync"
)

type FileManager struct {
	*sync.RWMutex
	files map[string]*FileInfo
}

func NewEmptyFileManager() *FileManager {
	return &FileManager{
		RWMutex: &sync.RWMutex{},
		files: map[string]*FileInfo{},
	}
}

func (fm *FileManager) fileInfoExists(name string) bool{
	fm.RLock()
	defer fm.RUnlock()
	_, exists := fm.files[name]
	return exists
}

func (fm *FileManager) AddFileInfo(name, hashString string, user User) error {
	hash := NewHash(hashString)
	if fm.fileInfoExists(name) {
		if fm.files[name].HasSameHash(hash) {
			err := fm.files[name].AddHolder(user)
			log.Printf("while adding %s as holder for %s: %v",user.name, name, err)
			return fmt.Errorf("you have already uploaded %s", name)
		} else {
			return fmt.Errorf("different file named %q already exists", name)
		}
	}
	fm.Lock()
	defer fm.Unlock()
	fm.files[name] = NewFileInfo(name, hash, map[string]User{user.name: user})
	return nil
}

func (fm *FileManager) DeleteFileInfo(name string) error {
	if fm.fileInfoExists(name){
		return fmt.Errorf("file info named %q does not exist", name)
	}
	fm.Lock()
	defer fm.Unlock()
	delete(fm.files, name)
	return nil
}

func (fm *FileManager) RemoveUserFromOwners(username string) error {
	fm.Lock()
	defer fm.Unlock()
	for i, file := range fm.files {
		if file.HasHolder(username){
			if err := fm.files[i].RemoveHolder(username); err != nil {
				return fmt.Errorf("while deleting holder from file %s: %w", file.name, err)
			}
		}
	}
	return nil
}
