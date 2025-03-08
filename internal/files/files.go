package files

import (
	"github.com/mat-sik/file-server-go/internal/envs"
	"os"
	"path/filepath"
)

func BuildClientFilePath(fileName string) string {
	return filepath.Join(envs.ClientStoragePath, fileName)
}

func buildServerFilePath(fileName string) string {
	return filepath.Join(envs.ServerStoragePath, fileName)
}

func getServerStoredFilenames() []string {
	return getAllFilenames(envs.ServerStoragePath)
}

func getAllFilenames(path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	filenames := make([]string, len(entries))
	for _, entry := range entries {
		filenames = append(filenames, entry.Name())
	}
	return filenames
}

func SizeOf(f *os.File) (int, error) {
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return int(stat.Size()), nil
}

func Close(f *os.File) {
	if err := f.Close(); err != nil {
		panic(err)
	}
}
