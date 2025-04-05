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
	rwMutex  *sync.RWMutex
	filename string
}

func (fh FileHandle) ExecuteReadOP(readOP func(string) error) error {
	fh.rwMutex.RLock()
	defer fh.rwMutex.RUnlock()
	return readOP(fh.filename)
}

func (fh FileHandle) ExecuteWriteOP(writeOP func(string) error) error {
	fh.rwMutex.Lock()
	defer fh.rwMutex.Unlock()
	return writeOP(fh.filename)
}

func NewFileHandle(filename string) FileHandle {
	return FileHandle{
		filename: filename,
		rwMutex:  &sync.RWMutex{},
	}
}

func (fh FileHandle) NewReadLockedFile() (*ReadLockedFile, error) {
	fh.rwMutex.RLock()

	file, err := os.Open(fh.filename)
	if err != nil {
		fh.rwMutex.RUnlock()
		return nil, err
	}

	return &ReadLockedFile{
		rwMutex: fh.rwMutex,
		file:    file,
	}, nil
}

type ReadLockedFile struct {
	rwMutex *sync.RWMutex
	file    *os.File
}

func (f *ReadLockedFile) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *ReadLockedFile) Size() (n int, err error) {
	return SizeOf(f.file)
}

func (f *ReadLockedFile) Close() error {
	defer f.rwMutex.RUnlock()
	return f.file.Close()
}
