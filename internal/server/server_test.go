package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mat-sik/file-server-go/internal/client"
	"github.com/mat-sik/file-server-go/internal/envs"
	"github.com/mat-sik/file-server-go/internal/files"
	"github.com/mat-sik/file-server-go/internal/message"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	setUpTestEnvs()
	setUpTestDirs()
	defer cleanTestDirs()

	m.Run()
}

func Test_shouldGetFileFromServerDeleteItOnServerAndPutItToServerUsingTheSameConnection(t *testing.T) {
	filename := "threeStepsTest.txt"
	serverFilePath := filepath.Join(testServerStoragePath, filename)
	createFile(serverFilePath, 1024*1024)

	// when
	cancel := runServerBlockTillListening()
	defer cancel()

	// and when
	webClient := getClient()

	getFileReq := message.GetFileRequest{Filename: filename}
	err := webClient.Run(getFileReq)

	// then
	if err != nil {
		t.Fatal(err)
	}
	clientFilePath := filepath.Join(testClientStoragePath, filename)
	if !filesEqual(clientFilePath, serverFilePath) {
		t.Fatalf("file not equal")
	}

	// and when
	delFileReq := message.DeleteFileRequest{Filename: filename}
	err = webClient.Run(delFileReq)

	// then
	if err != nil {
		t.Fatal(err)
	}
	if fileExists(serverFilePath) {
		t.Fatalf("file exists, but should have been deleted")
	}

	// and when
	putFileReq := message.PutFileRequest{Filename: filename}
	err = webClient.Run(putFileReq)

	// then
	if err != nil {
		t.Fatal(err)
	}
	if !filesEqual(serverFilePath, clientFilePath) {
		t.Fatalf("file not equal")
	}
}

func Test_shouldGetFileFromServer(t *testing.T) {
	// given
	filename := "getFileTest.txt"
	serverFilePath := filepath.Join(testServerStoragePath, filename)
	createFile(serverFilePath, 1024*1024)

	// when
	cancel := runServerBlockTillListening()
	defer cancel()

	// and when
	webClient := getClient()

	getFileReq := message.GetFileRequest{Filename: filename}
	err := webClient.Run(getFileReq)

	// then
	if err != nil {
		t.Fatal(err)
	}
	clientFilePath := filepath.Join(testClientStoragePath, filename)
	if !filesEqual(clientFilePath, serverFilePath) {
		t.Fatalf("file not equal")
	}
}

func Test_shouldPutFileToServer(t *testing.T) {
	// given
	filename := "putFileTest.txt"
	clientFilePath := filepath.Join(testClientStoragePath, filename)
	createFile(clientFilePath, 1024*1024)

	// when
	cancel := runServerBlockTillListening()
	defer cancel()

	// and when
	webClient := getClient()

	putFileReq := message.PutFileRequest{Filename: filename}
	err := webClient.Run(putFileReq)

	// then
	if err != nil {
		t.Fatal(err)
	}
	serverFilePath := filepath.Join(testServerStoragePath, filename)
	if !filesEqual(serverFilePath, clientFilePath) {
		t.Fatalf("file not equal")
	}
}

func Test_shouldDeleteFileFromServer(t *testing.T) {
	// given
	filename := "deleteFileTest.txt"
	serverFilePath := filepath.Join(testServerStoragePath, filename)
	createFile(serverFilePath, 1024*1024)

	// when
	cancel := runServerBlockTillListening()
	defer cancel()

	// and when
	webClient := getClient()

	delFileReq := message.DeleteFileRequest{Filename: filename}
	err := webClient.Run(delFileReq)

	// then
	if err != nil {
		t.Fatal(err)
	}
	if fileExists(serverFilePath) {
		t.Fatalf("file exists, but should have been deleted")
	}
}

func runServerBlockTillListening() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go runServer(ctx, wg)
	wg.Wait()

	return cancel
}

func runServer(ctx context.Context, wg *sync.WaitGroup) {
	addr := fmt.Sprintf(":%d", port)

	if err := runWithWaitGroup(ctx, wg, addr); err != nil {
		panic(err)
	}
}

func getClient() client.Client {
	addr := fmt.Sprintf(":%d", port)

	webClient, err := client.NewClient(addr)
	if err != nil {
		panic(err)
	}

	return webClient
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}

func filesEqual(firstPath string, secondPath string) bool {
	firstFile, err := os.Open(firstPath)
	if err != nil {
		panic(err)
	}
	defer files.LoggedClose(firstFile)
	secondFile, err := os.Open(secondPath)
	if err != nil {
		panic(err)
	}
	defer files.LoggedClose(secondFile)

	return fileLengthEqual(firstFile, secondFile) && fileContentsEqual(firstFile, secondFile)
}

func fileContentsEqual(firstFile *os.File, secondFile *os.File) bool {
	firstBuffer := copyBuffer[:len(copyBuffer)/2]
	secondBuffer := copyBuffer[len(copyBuffer)/2:]

	for {
		firstN, firstErr := firstFile.Read(firstBuffer)
		secondN, secondErr := secondFile.Read(secondBuffer)
		if (firstErr != nil && secondErr == nil) || (firstErr == nil && secondErr != nil) {
			return false
		} else if firstErr != nil && secondErr != nil {
			return true
		}
		if firstN != secondN {
			return false
		}
		if !bytes.Equal(firstBuffer[:firstN], secondBuffer[:secondN]) {
			return false
		}
	}
}

func fileLengthEqual(firstFile *os.File, secondFile *os.File) bool {
	firstSize, err := files.SizeOf(firstFile)
	if err != nil {
		panic(err)
	}
	secondSize, err := files.SizeOf(secondFile)
	if err != nil {
		panic(err)
	}

	return firstSize == secondSize
}

func createFile(path string, size int) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer files.LoggedClose(file)

	for i := 0; i < min(size, len(copyBuffer)); i++ {
		copyBuffer[i] = 'x'
	}

	bytesWritten := 0
	for bytesWritten < size {
		remaining := size - bytesWritten

		writeSize := min(remaining, len(copyBuffer))
		n, err := file.Write(copyBuffer[:writeSize])
		if err != nil {
			panic(err)
		}

		bytesWritten += n
	}

	slog.Info("Created file", "path", path, "size", bytesWritten)
}

func setUpTestDirs() {
	dirs := []string{testServerStoragePath, testClientStoragePath}
	createDirs(dirs)
}

func createDirs(dirs []string) {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			cleanTestDirs()
			panic(err)
		}
	}
}

func cleanTestDirs() {
	if err := os.RemoveAll(pathToTest); err != nil {
		panic(err)
	}
}

func setUpTestEnvs() {
	envs.ClientStoragePath = testClientStoragePath
	envs.ServerStoragePath = testServerStoragePath
}

var (
	pathToRoot            = "../../"
	pathToTest            = filepath.Join(pathToRoot, "test")
	testServerStoragePath = filepath.Join(pathToTest, "server/storage")
	testClientStoragePath = filepath.Join(pathToTest, "client/storage")

	port = 33303

	copyBuffer = make([]byte, 64*1024)
)
