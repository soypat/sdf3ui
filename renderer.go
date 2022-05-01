package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"github.com/fsnotify/fsnotify"
	"github.com/soypat/sdf/render"
)

type rendererHandler struct {
	mu            sync.Mutex
	w             *fsnotify.Watcher
	watchCtx      context.Context
	server        shape3DServer
	refreshPeriod time.Duration
}

func newRendererHandler() *rendererHandler {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	return &rendererHandler{
		w:             w,
		refreshPeriod: time.Second,
	}
}

func (r *rendererHandler) SetFileTarget(filename string) error {
	if len(r.w.WatchList()) > 0 {
		r.w.Remove(r.w.WatchList()[0])
	}
	return r.w.Add(filename)
}

func (r *rendererHandler) Start(ctx context.Context) error {
	for {
		lastCmd := time.Now()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case event, ok := <-r.w.Events:
				if !ok {
					return errors.New("watcher Events closed.")
				}

				if event.Op&fsnotify.Write == fsnotify.Write && time.Since(lastCmd) > r.refreshPeriod {
					lastCmd = time.Now()
					log.Println("event:", event)
					log.Println("The modified file:", event.Name)
					err := r.renderFile(ctx, event.Name)
					if err != nil {
						log.Println("error running command:", err)
					}
				}

			case err, ok := <-r.w.Errors:
				if !ok {
					return errors.New("watcher.Errors channel closed")
				}
				log.Println("watcher error:", err)
			}
		}
	}
}

func (r *rendererHandler) renderFile(ctx context.Context, filename string) error {
	output, err := exec.Command("go", "run", filename, "stdout").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s\n%s", string(output), err.Error())
	}
	lenTri := len(output) / int(unsafe.Sizeof(render.Triangle3{}))
	sli := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&output[0])),
		Len:  lenTri,
		Cap:  lenTri,
	}
	triSlice := *(*[]render.Triangle3)(unsafe.Pointer(&sli))
	r.server.SetShape(triSlice)
	return nil
}
