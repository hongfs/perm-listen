package main

import (
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
)

// LISTEN_PATH=/home/www/ LISTEN_EXTENSION=.log LISTEN_USER=hongfs go run main.go

func init() {
	file, err := os.OpenFile("/tmp/perm-listen.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)

	if err != nil {
		panic(err)
	}

	log.SetOutput(file)
	log.SetPrefix("[listen]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

func main() {
	err := handle()

	if err != nil {
		log.Fatal(err)
	}
}

func handle() error {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return err
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					addListen(watcher, event.Name)
					continue
				}

				if !event.Has(fsnotify.Chmod) {
					continue
				}

				err := chmodFile(event.Name)

				if err != nil {
					log.Printf("chmod file error: %s\n", err)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("error:", err)
			}
		}
	}()

	err = addListen(watcher, os.Getenv("LISTEN_PATH"))

	if err != nil {
		return err
	}

	<-make(chan struct{})

	return nil
}

var listenMap = new(sync.Map)

func addListen(w *fsnotify.Watcher, path string) error {
	if !fileutil.IsDir(path) {
		return nil
	}

	if _, ok := listenMap.Load(path); ok {
		return nil
	}

	log.Printf("add listen: %s\n", path)

	err := w.Add(path)

	if err != nil {
		return err
	}

	listenMap.Store(path, true)

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		return addListen(w, path)
	})
}

func chmodFile(name string) error {
	log.Println("modified file:", name)

	extension := filepath.Ext(name)

	if extension != os.Getenv("LISTEN_EXTENSION") {
		return nil
	}

	info, err := os.Stat(name)

	if err != nil {
		return err
	}

	stat, ok := info.Sys().(*syscall.Stat_t)

	if !ok {
		return errors.New("get file stat error")
	}

	u, err := user.Lookup(os.Getenv("LISTEN_USER"))

	if err != nil {
		return err
	}

	if fmt.Sprintf("%d", stat.Uid) == u.Uid {
		return nil
	}

	uid, err := strconv.ParseInt(u.Uid, 10, 64)

	if err != nil {
		return err
	}

	return os.Chown(name, int(uid), int(stat.Gid))
}
