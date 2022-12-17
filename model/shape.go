package model

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"time"

	"gonum.org/v1/gonum/spatial/r3"
)

const (
	HTTPHost    = "http://[::]"
	HTTPAddr    = ":8080"
	HTTPBaseURL = HTTPHost + HTTPAddr
	// Websocket sub protocol.
	WSSubprotocol   = "sdf3ui"
	ShapeEndpoint   = "/" + WSSubprotocol + "/getShape"
	SaveSTLEndpoint = "/createSTL"
	WSTimeout       = 60 * time.Minute
	kilobyte        = 1000
	megabyte        = 1000 * kilobyte
	WSReadLimit     = 32 * kilobyte
)

// Shape3D contains 3D shape information.
type Shape3D struct {
	// Ctx stores shape status. This context is cancelled when
	// a newer shape is available. The previous shape is said
	// to be stale in this case.
	ctx context.Context
	// Triangles has render data.
	Triangles []r3.Triangle
	// Sequence number of shape.
	Seq uint
}

func (s Shape3D) String() string {
	if s.Context().Err() != nil {
		return "stale shape"
	}
	return fmt.Sprintf("seq:%d, faces:%d", s.Seq, len(s.Triangles))
}

func (s *Shape3D) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *Shape3D) Context() context.Context {
	if s.ctx == nil {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		return ctx
	}
	return s.ctx
}

type ShapeStatus string

const (
	ShapeStale ShapeStatus = "shape is stale"
	ShapeOK    ShapeStatus = "shape is up to date"
)

type ServerStatus struct {
	ShapeSeq  uint
	ShapeName string
}

func (s *Shape3D) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(s)
}

func (s *Shape3D) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(s)
}
