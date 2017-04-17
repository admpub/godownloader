package iotools

import (
	"os"
	"sync"
)

type SafeFile struct {
	*os.File
	lock     sync.Mutex
	filePath string
	closed   bool
}

func (sf *SafeFile) WriteAt(b []byte, off int64) (n int, err error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	return sf.File.WriteAt(b, off)
}

func (sf *SafeFile) Sync() error {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	return sf.File.Sync()
}

func (sf *SafeFile) Close() error {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	if sf.closed {
		return nil
	}
	sf.closed = true
	return sf.File.Close()
}

func (sf *SafeFile) ReOpen() error {
	if !sf.closed {
		return nil
	}
	sf.lock.Lock()
	defer sf.lock.Unlock()
	f, err := os.OpenFile(sf.filePath, os.O_RDWR, 0666)
	sf.File = f
	return err
}

func OpenSafeFile(name string) (file *SafeFile, err error) {
	f, err := os.OpenFile(name, os.O_RDWR, 0666)
	return &SafeFile{File: f, filePath: name}, err
}

func CreateSafeFile(name string) (file *SafeFile, err error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	return &SafeFile{File: f, filePath: name}, err
}
