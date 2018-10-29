package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"github.com/edsrzf/mmap-go"
	"syscall"
	"time"
)

func main() {
	// filelock()
	memoryMap()
}

func memoryMap() {
	var testData = []byte("0123456789ABCDEF")
	var testPath = filepath.Join(os.TempDir(), "testdata")
	err := ioutil.WriteFile(testPath, testData, 0644)
	e(err)

	f, err := os.OpenFile(testPath, os.O_RDWR, 0644)
	e(err)
	defer f.Close()
	m, err := mmap.Map(f, mmap.RDWR, 0644)
	e(err)
	defer m.Unmap()

	m[9] = 'X'
	m.Flush()

	fileData, err := ioutil.ReadAll(f)
	e(err)

	fmt.Printf("original: %s\n", testData)
  fmt.Printf("mmap: %s\n", m)
  fmt.Printf("file: %s\n", fileData)
}

type FileLock struct {
	l  sync.Mutex
	fd int
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}

func NewFileLock(filename string) *FileLock {
	fd, err := syscall.Open(filename, syscall.O_CREAT|syscall.O_RDONLY, 0750)
	e(err)
	return &FileLock{fd: fd}
}

func (m *FileLock) Lock() {
	m.l.Lock()
	if err := syscall.Flock(m.fd, syscall.LOCK_EX); err != nil {
		panic(err)
	}
}

func (m *FileLock) Unlock() {
	if err := syscall.Flock(m.fd, syscall.LOCK_UN); err != nil {
		panic(err)
	}
	m.l.Unlock()
}

func filelock() {
	fmt.Println("==in_first==")
	fl := NewFileLock("hoge")
	fl.Lock()
	fmt.Println("==first lock done")
	time.Sleep(10 * time.Second)
	fmt.Println("==first sleep done")
	fl.Unlock()
	fmt.Println("==first unlock done")
}
