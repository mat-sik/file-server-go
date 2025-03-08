package envs

import "os"

var (
	ClientStoragePath = clientStoragePath()
	ServerStoragePath = serverStoragePath()
)

func serverStoragePath() string {
	return os.Getenv("SERVER_STORAGE_PATH")
}

func clientStoragePath() string {
	return os.Getenv("CLIENT_STORAGE_PATH")
}
