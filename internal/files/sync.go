package files

import (
	"os"
	"path/filepath"
	"sync"
)

type SyncService struct {
	files *sync.Map
}

func (s *SyncService) AddFile(filename string) FileHandle {
	path := buildServerFilePath(filename)
	fileHandler := NewFileHandle(path)
	s.files.Store(path, fileHandler)
	return fileHandler
}

func (s *SyncService) GetFile(filename string) (FileHandle, bool) {
	path := buildServerFilePath(filename)
	fileHandle, ok := s.files.Load(path)
	if !ok {
		return FileHandle{}, false
	}
	return fileHandle.(FileHandle), true
}

func (s *SyncService) RemoveFile(filename string) error {
	path := buildServerFilePath(filename)
	value, ok := s.files.Load(path)
	if !ok {
		return os.ErrNotExist
	}
	fileHandle := value.(FileHandle)
	if err := fileHandle.ExecuteWriteOP(os.Remove); err != nil {
		return err
	}
	s.files.Delete(path)
	return nil
}

func (s *SyncService) GetAllFilenames() []string {
	var filenames []string
	s.files.Range(func(key, value interface{}) bool {
		path := key.(string)
		filenames = append(filenames, filepath.Base(path))
		return true
	})
	return filenames
}

func NewService() SyncService {
	fileService := SyncService{
		files: &sync.Map{},
	}

	filenames := getServerStoredFilenames()
	for _, filename := range filenames {
		fileService.AddFile(filename)
	}

	return fileService
}

type FileHandle struct {
	*sync.RWMutex
	filename string
}

func (fh FileHandle) ExecuteReadOP(readOP func(string) error) error {
	fh.RLock()
	defer fh.RUnlock()
	return readOP(fh.filename)
}

func (fh FileHandle) ExecuteWriteOP(writeOP func(string) error) error {
	fh.Lock()
	defer fh.Unlock()
	return writeOP(fh.filename)
}

func NewFileHandle(filename string) FileHandle {
	return FileHandle{
		filename: filename,
		RWMutex:  &sync.RWMutex{},
	}
}

func (fh FileHandle) NewReadLockedFile() (ReadLockedFile, error) {
	fh.RLock()

	file, err := os.Open(fh.filename)
	if err != nil {
		fh.RUnlock()
		return ReadLockedFile{}, err
	}

	return ReadLockedFile{
		RWMutex: fh.RWMutex,
		File:    file,
	}, nil
}

type ReadLockedFile struct {
	*sync.RWMutex
	*os.File
}

func (f *ReadLockedFile) Close() error {
	defer f.RUnlock()
	return f.File.Close()
}
