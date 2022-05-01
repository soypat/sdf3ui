package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
	"vecty-templater-project/model"

	"github.com/soypat/sdf/render"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type shape3DServer struct {
	mu    sync.Mutex
	shape model.Shape3D
	// staleFunc is a function that is called when shape
	// becomes stale. It cancels the Shape3D's context.
	staleFunc func()
}

// ServeHTTP is a basic websocket implementation for reading/writing a TODO list
// from a websocket.
func (s *shape3DServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:       []string{model.WSSubprotocol},
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	c.SetReadLimit(model.MaxRenderSize)
	if c.Subprotocol() != model.WSSubprotocol {
		c.Close(websocket.StatusPolicyViolation, "client must speak the "+model.WSSubprotocol+" subprotocol")
		return
	}
	log.Println("websocket connection established")
	for {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
		err = s.sendShape(ctx, c)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		cancel()
		if err != nil {
			log.Printf("failed to echo with %v: %v\n", r.RemoteAddr, err)
			return
		}
	}
}

// sendShape sends 3d model data over websocket connection. It fails to send
// if shape becomes stale.
func (t *shape3DServer) sendShape(ctx context.Context, c *websocket.Conn) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.shape.Ctx == nil {
		t.shape.Ctx = ctx // temp workaround
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*200)
	defer cancel()
	go func() {
		// This goroutine cancels sendShape if shape becomes stale.
		select {
		case <-ctx.Done():
		case <-t.shape.Ctx.Done():
		}
		cancel()
	}()
	return wsjson.Write(ctx, c, t.shape)
}

// SetShape sets the render data and handles Shape3D context
// and sequence number.
func (t *shape3DServer) SetShape(data []render.Triangle3) {
	if t.staleFunc != nil {
		t.staleFunc()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	t.staleFunc = cancel
	seq := t.shape.Seq + 1
	t.shape = model.Shape3D{
		Ctx:       ctx,
		Triangles: data,
		Seq:       seq,
	}
}
