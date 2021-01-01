package memory

import (
	"fmt"
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
			return fmt.Errorf("while adding %s as holder for %s: %v", user.Name, name, err)
		} else {
			return fmt.Errorf("different file named %q already exists", name)
		}
	}
	fm.Lock()
	defer fm.Unlock()
	fm.files[name] = NewFileInfo(name, hash, map[string]User{user.Name: user})
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

func (fm *FileManager) GetFileInfo(name string) (*FileInfo, error) {
	if fm.fileInfoExists(name){
		return nil, fmt.Errorf("file info named %q does not exist", name)
	}
	//TODO: Other goroutine could delete meanwhile
	fm.Lock()
	defer fm.Unlock()
	return fm.files[name], nil
}

func (fm *FileManager) DeleteUserFromFileInfo(filename string, user User) error {
	if !fm.fileInfoExists(filename) {
			return fmt.Errorf("file named %q does not exist", filename)
	}

	if err := fm.files[filename].RemoveHolder(user.Name); err != nil {
		return fmt.Errorf("while removing %s as holder for %s: %v", user.Name, filename, err)
	}

	if !fm.files[filename].HasAnyHolders(){
		fm.Lock()
		defer fm.Unlock()
		if err := fm.DeleteFileInfo(filename); err != nil {
			return fmt.Errorf("while removing %s as holder for %s: %v", user.Name, filename, err)
		}
	}

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
