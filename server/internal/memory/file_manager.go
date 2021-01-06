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
		files:   map[string]*FileInfo{},
	}
}

func (fm *FileManager) AddFileInfo(name, hashString string, user User) error {
	hash := NewHash(hashString)
	if exists, err := fm.safeCheckForExistence(name, hash); exists {
		if err := fm.files[name].AddHolder(user); err != nil {
			return fmt.Errorf("while adding %q as holder for %q: %v", user.Name, name, err)
		}
	} else if err != nil {
		return err
	}

	fm.Lock()
	defer fm.Unlock()
	fm.files[name] = NewFileInfo(name, hash, map[string]User{user.Name: user})
	return nil
}

func (fm *FileManager) DeleteFileInfo(name string) error {
	fm.Lock()
	defer fm.Unlock()
	delete(fm.files, name)
	return nil
}

func (fm *FileManager) GetFileInfo(name string) (*FileInfo, error) {
	fm.RLock()
	defer fm.RUnlock()
	if !fm.fileInfoExists(name) {
		return nil, fmt.Errorf("file info named %q does not exist", name)
	}
	return fm.files[name], nil
}

func (fm *FileManager) DeleteUserFromFileInfo(filename string, user User) error {
	fm.RLock()
	if !fm.fileInfoExists(filename) {
		fm.RUnlock()
		return fmt.Errorf("file named %q does not exist", filename)
	}
	fm.RUnlock()

	if err := fm.files[filename].RemoveHolder(user.Name); err != nil {
		return fmt.Errorf("while removing %s as holder for %s: %v", user.Name, filename, err)
	}

	if !fm.files[filename].HasAnyHolders() {
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
		if err := fm.files[i].RemoveHolder(username); err != nil {
			return fmt.Errorf("while deleting holder from file %s: %w", file.name, err)
		}
	}
	return nil
}

func (fm *FileManager) fileInfoExists(name string) bool {
	_, exists := fm.files[name]
	return exists
}

type FileAlreadyExistsError struct {
	filename string
}

func (f FileAlreadyExistsError) Error() string {
	return fmt.Sprintf("Another file named %q already exists", f.filename)
}

func (fm *FileManager) safeCheckForExistence(name string, hash Hash) (bool, error) {
	fm.RLock()
	defer fm.RUnlock()
	if !fm.fileInfoExists(name) {
		return false, nil
	}
	if !fm.files[name].HasSameHash(hash) {
		return false, FileAlreadyExistsError{}
	}
	return true, nil
}
