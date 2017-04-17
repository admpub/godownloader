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
	opened   []*os.File
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
	for i, f := range sf.opened {
		if f == sf.File {
			start := i + 1
			if start < len(sf.opened) {
				sf.opened = append(sf.opened[0:i], sf.opened[start:0]...)
			} else {
				sf.opened = sf.opened[0:i]
			}
		}
	}
	return sf.File.Close()
}

func (sf *SafeFile) Opened() int {
	return len(sf.opened)
}

func (sf *SafeFile) OpenedFiles() []*os.File {
	return sf.opened
}

func (sf *SafeFile) CloseAll() {
	for _, f := range sf.opened {
		f.Close()
	}
	sf.opened = []*os.File{}
}

func (sf *SafeFile) ReOpen() error {
	sf.lock.Lock()
	defer sf.lock.Unlock()
	if !sf.closed {
		return nil
	}
	f, err := os.OpenFile(sf.filePath, os.O_RDWR, 0666)
	sf.File = f
	sf.opened = append(sf.opened, f)
	return err
}

func OpenSafeFile(name string) (file *SafeFile, err error) {
	file = &SafeFile{filePath: name, opened: []*os.File{}}
	err = file.ReOpen()
	return
}

func CreateSafeFile(name string) (file *SafeFile, err error) {
	var f *os.File
	f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	file = &SafeFile{File: f, filePath: name, opened: []*os.File{}}
	if err == nil {
		file.opened = append(file.opened, f)
	}
	return
}
