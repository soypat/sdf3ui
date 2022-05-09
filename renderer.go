package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/crypto/md4"

	"github.com/soypat/sdf3ui/uirender"

	"github.com/fsnotify/fsnotify"
)

type rendererHandler struct {
	mu            sync.Mutex
	w             *fsnotify.Watcher
	watchCtx      context.Context
	server        shape3DServer
	refreshPeriod time.Duration
	// make sure only rendering changed shapes.
	outputHash uint64
}

func newRendererHandler() *rendererHandler {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	return &rendererHandler{
		w:             w,
		refreshPeriod: 1000 * time.Millisecond,
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
				if event.Op&fsnotify.Write == fsnotify.Write {
					elapsed := time.Since(lastCmd)
					if elapsed < r.refreshPeriod {
						log.Printf("too soon to rerender %s < %s\n", elapsed, r.refreshPeriod)
						continue
					}
					lastCmd = time.Now()
					err := r.renderFile(ctx, event.Name)
					if err != nil {
						log.Println(err)
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

func (r *rendererHandler) renderFile(ctx context.Context, filename string) (err error) {
	var outHash uint64
	defer func() {
		// log.Printf("had hash %x got hash %x (matching hashes means render success)", r.outputHash, outHash)
	}()

	log.Println("[REND] go run", filename)
	output, err := exec.Command("go", "run", filename, "stdout").CombinedOutput()
	if err != nil {
		return fmt.Errorf("[ERRO] in command:\n%s\n%s", string(output), err.Error())
	}
	outHash, err = hashReader(bytes.NewReader(output))
	if outHash == r.outputHash {
		return fmt.Errorf("[SKIP] file output has not changed since last render")
	}

	triangles, err := uirender.DecodeAll(bytes.NewReader(output))
	if err != nil {
		return fmt.Errorf("[ERRO] decoding: %s", err.Error())
	}
	r.outputHash = outHash
	r.server.SetShape(filename, triangles)
	return nil
}

func hashReader(r io.Reader) (uint64, error) {
	hsh := md4.New()
	_, err := io.Copy(hsh, r)
	if err != nil {
		return 0, err
	}
	sum := hsh.Sum(nil)
	v1 := binary.BigEndian.Uint64(sum[:8])
	// v2 := binary.BigEndian.Uint64(sum[8:])
	return v1, nil
}
