// Copyright © 2017 Kaleo Cheng <kaleocheng@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package watcher

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type RecursiveWatcher struct {
	*fsnotify.Watcher
	interval time.Duration
	Files    chan []string
}

func NewRecursiveWatcher(path string, interval time.Duration) (*RecursiveWatcher, error) {
	folders := Subfolders(path)
	if len(folders) == 0 {
		return nil, errors.New("no folders to watch")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	rw := &RecursiveWatcher{
		Watcher:  watcher,
		interval: interval,
		Files:    make(chan []string, 10),
	}
	for _, folder := range folders {
		rw.AddFolder(folder)
	}
	return rw, nil
}

func (watcher *RecursiveWatcher) AddFolder(folder string) {
	err := watcher.Add(folder)
	if err != nil {
		log.Println("Error watching: ", folder, err)
	}
}

func (watcher *RecursiveWatcher) Run() {
	go func() {
		tick := time.Tick(watcher.interval)
		names := make([]string, 0)
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					fi, err := os.Stat(event.Name)
					if err != nil {
						log.Println(err)
						break
					}
					if fi.IsDir() {
						if !shouldIgnoreFile(filepath.Base(event.Name)) {
							watcher.AddFolder(event.Name)
						}
					}
				}

				if event.Op&fsnotify.Write == fsnotify.Write {

					if !shouldIgnoreFile(event.Name) {
						names = append(names, event.Name)
					}
				}

			case <-tick:
				if len(names) == 0 {
					continue
				}
				watcher.Files <- names
				names = make([]string, 0)
			case err := <-watcher.Errors:
				log.Println("error ", err)
			}
		}
	}()
}

func Subfolders(path string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			if shouldIgnoreFile(name) && name != "." && name != ".." {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

func shouldIgnoreFile(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") || name == "kanna"
}
