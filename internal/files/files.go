package files

import (
	"github.com/mat-sik/file-server-go/internal/envs"
	"os"
	"path/filepath"
)

func BuildServerFilePath(fileName string) string {
	return filepath.Join(envs.ServerDBPath, fileName)
}

func BuildClientFilePath(fileName string) string {
	return filepath.Join(envs.ClientDBPath, fileName)
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
