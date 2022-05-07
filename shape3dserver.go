package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
	"vecty-templater-project/model"

	"github.com/soypat/sdf/render"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type shape3DServer struct {
	mu    sync.Mutex
	shape model.Shape3D
	// staleFunc is a function that is called when shape
	// becomes stale. It cancels the Shape3D's context.
	staleFunc   func()
	shapeNotify chan struct{}
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
	l := rate.NewLimiter(rate.Every(500*time.Millisecond), 2)
	log.Println("websocket connection established")

	if s.shapeNotify == nil {
		s.shapeNotify = make(chan struct{})
	}
	for range s.shapeNotify {
		log.Println("got shapeNotify")
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		err = s.sendStatus(ctx, c, l)
		cancel()
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			log.Printf("failed to echo with %v: %v\n", r.RemoteAddr, err)
			return
		}
		time.Sleep(500 * time.Millisecond) // burst rate limiting
	}
}

// sendStatus sends 3d model data over websocket connection. It fails to send
// if shape becomes stale.
func (t *shape3DServer) sendStatus(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.shape.Context().Err() != nil {
		t.shape.SetContext(ctx) // temp workaround
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		// This goroutine cancels sendShape if shape becomes stale.
		select {
		case <-ctx.Done():
		case <-t.shape.Context().Done():
		}
		cancel()
	}()
	status := model.ServerStatus{
		ShapeSeq: t.shape.Seq,
	}
	err := l.Wait(ctx)
	if err != nil {
		return err
	}
	return wsjson.Write(ctx, c, status)
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
		Triangles: data,
		Seq:       seq,
	}
	t.shape.SetContext(ctx)
busysend:
	for {
		select {
		case t.shapeNotify <- struct{}{}:
			// send until no more requests waiting.
		default:
			break busysend
		}
	}

}

func (t *shape3DServer) serveShapeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	log.Println("serveShapeHTTP request")
	t.mu.Lock()
	defer t.mu.Unlock()
	log.Println("encoding shape")
	err := t.shape.Encode(w)
	if err != nil {
		log.Println("error encoding shape", err)
		return
	}
	log.Println("shape encode success")
}
