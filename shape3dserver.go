package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/soypat/sdf3ui/model"

	"github.com/soypat/sdf/render"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type shape3DServer struct {
	mu    sync.Mutex
	shape model.Shape3D
	stat  model.ServerStatus
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

	c.SetReadLimit(model.WSReadLimit)
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
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = s.sendStatus(ctx, c, l)
		cancel()
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			c.Close(websocket.StatusInternalError, "normal closure")
			return
		}
		if err != nil {
			c.Close(websocket.StatusInternalError, err.Error())
			log.Printf("failed to send status to %v: %v\n", r.RemoteAddr, err)
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
			log.Println("shape became stale during status send")
		}
		cancel()
	}()
	err := l.Wait(ctx)
	if err != nil {
		return err
	}
	return wsjson.Write(ctx, c, t.stat)
}

// SetShape sets the render data and handles Shape3D context
// and sequence number.
func (t *shape3DServer) SetShape(name string, data []render.Triangle3) {
	if t.staleFunc != nil {
		t.staleFunc()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	t.staleFunc = cancel
	seq := t.shape.Seq + 1
	t.stat = model.ServerStatus{
		ShapeSeq:  seq,
		ShapeName: name,
	}
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
	defer log.Printf("request to encode %s ended\n", t.shape.String())
	w.Header().Add("Content-Type", "application/octet-stream")
	w.WriteHeader(200)
	t.mu.Lock()
	defer t.mu.Unlock()
	err := t.shape.Encode(w)
	if err != nil {
		log.Printf("[ERR] encoding shape %s: %s", t.shape.String(), err)
		return
	}
}

func (t *shape3DServer) createSTLHandler(w http.ResponseWriter, req *http.Request) {
	t.mu.Lock()
	defer t.mu.Unlock()
	filename := req.URL.Query().Get("name")
	fp, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = render.WriteSTL(fp, t.shape.Triangles)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(200)
}
