package files

import (
	"os"
)

func GetSize(f *os.File) (int, error) {
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
