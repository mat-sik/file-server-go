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

func Test_shouldReturnedAllStoredFilenames(t *testing.T) {
	// given
	serverFilename1 := "serverFilenameA"
	serverPath1 := filepath.Join(testServerStoragePath, serverFilename1)
	createFile(serverPath1, 1024*1024)

	serverFilename2 := "serverFilenameAA"
	serverPath2 := filepath.Join(testServerStoragePath, serverFilename2)
	createFile(serverPath2, 1024*1024)

	serverFilename3 := "serverFilenameBB"
	serverPath3 := filepath.Join(testServerStoragePath, serverFilename3)
	createFile(serverPath3, 1024*1024)

	serverFilename4 := "serverFilenameAC"
	serverPath4 := filepath.Join(testServerStoragePath, serverFilename4)
	createFile(serverPath4, 1024*1024)

	// when
	cancel := runServerBlockTillListening()
	defer cancel()

	// and when
	webClient := getClient()

	getFilenamesReq := message.GetFilenamesRequest{MatchRegex: ".*A.*"}
	err := webClient.Run(getFilenamesReq)

	// then
	if err != nil {
		t.Fatal(err)
	}
}

func Test_shouldPassRaceConditionTest(t *testing.T) {
	// given
	serverFilename1 := "serverFilename1"
	serverPath1 := filepath.Join(testServerStoragePath, serverFilename1)
	createFile(serverPath1, 1024*1024)

	serverFilename2 := "serverFilename2"
	serverPath2 := filepath.Join(testServerStoragePath, serverFilename2)
	createFile(serverPath2, 1024*1024)

	clientFilename1 := "clientFilename1"
	clientPath1 := filepath.Join(testClientStoragePath, clientFilename1)
	createFile(clientPath1, 1024*1024)

	clientFilename2 := "clientFilename2"
	clientPath2 := filepath.Join(testClientStoragePath, clientFilename2)
	createFile(clientPath2, 1024*1024)

	// when
	cancel := runServerBlockTillListening()
	defer cancel()

	// and when
	wg := &sync.WaitGroup{}
	wg.Add(4)
	go runRequest(wg, func(webClient client.Client) error {
		req := message.GetFileRequest{Filename: serverFilename1}
		for i := 0; i < 5; i++ {
			if err := webClient.Run(req); err != nil {
				return err
			}
		}
		return nil
	})
	go runRequest(wg, func(webClient client.Client) error {
		req := message.DeleteFileRequest{Filename: serverFilename1}
		if err := webClient.Run(req); err != nil {
			return err
		}
		return nil
	})
	go runRequest(wg, func(webClient client.Client) error {
		getReq := message.GetFileRequest{Filename: serverFilename2}
		for i := 0; i < 5; i++ {
			if err := webClient.Run(getReq); err != nil {
				return err
			}
		}
		putReq := message.PutFileRequest{Filename: clientFilename1}
		if err := webClient.Run(putReq); err != nil {
			return err
		}
		putReq = message.PutFileRequest{Filename: clientFilename2}
		if err := webClient.Run(putReq); err != nil {
			return err
		}
		return nil
	})
	go runRequest(wg, func(webClient client.Client) error {
		getReq := message.GetFileRequest{Filename: serverFilename2}
		for i := 0; i < 5; i++ {
			if err := webClient.Run(getReq); err != nil {
				return err
			}
		}
		for i := 0; i < 5; i++ {
			delReq := message.DeleteFileRequest{Filename: clientFilename1}
			if err := webClient.Run(delReq); err != nil {
				return err
			}
			delReq = message.DeleteFileRequest{Filename: clientFilename2}
			if err := webClient.Run(delReq); err != nil {
				return err
			}
		}
		return nil
	})

	wg.Wait()
}

func runRequest(wg *sync.WaitGroup, execRequest func(webClient client.Client) error) {
	webClient := getClient()
	if err := execRequest(webClient); err != nil {
		panic(err)
	}
	wg.Done()
}

func Test_shouldGetFileFromServerDeleteItOnServerAndPutItToServerUsingTheSameConnection(t *testing.T) {
	// given
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
