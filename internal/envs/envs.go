package envs

import "os"

var (
	ClientDBPath = getClientDBPath()
	ServerDBPath = getServerDBPath()
)

func getServerDBPath() string {
	return os.Getenv("SERVER_DB_PATH")
}

func getClientDBPath() string {
	return os.Getenv("CLIENT_DB_PATH")
}
