package envs

import "os"

var (
	ClientDBPath = clientDBPath()
	ServerDBPath = serverDBPath()
)

func serverDBPath() string {
	return os.Getenv("SERVER_DB_PATH")
}

func clientDBPath() string {
	return os.Getenv("CLIENT_DB_PATH")
}
