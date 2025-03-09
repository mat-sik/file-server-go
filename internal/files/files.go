package files

import (
	"github.com/mat-sik/file-server-go/internal/envs"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func BuildClientFilePath(filename string) string {
	return filepath.Join(envs.ClientStoragePath, filename)
}

func buildServerFilePath(filename string) string {
	return filepath.Join(envs.ServerStoragePath, filename)
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

func LoggedClose(f io.Closer) {
	if err := f.Close(); err != nil {
		slog.Error(err.Error())
	}
}
