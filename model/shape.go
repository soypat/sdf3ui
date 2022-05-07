package model

import (
	"context"
	"encoding/gob"
	"io"

	"github.com/soypat/sdf/render"
)

const (
	HTTPHost    = "http://[::]"
	HTTPAddr    = ":8080"
	HTTPBaseURL = HTTPHost + HTTPAddr
	// Websocket sub protocol.
	WSSubprotocol   = "sdf3ui"
	ShapeEndpoint   = "/" + WSSubprotocol + "/getShape"
	SaveSTLEndpoint = "/createSTL"
	megabyte        = 1000 * 1000
	MaxRenderSize   = 30 * megabyte
)

// Shape3D contains 3D shape information.
type Shape3D struct {
	// Ctx stores shape status. This context is cancelled when
	// a newer shape is available. The previous shape is said
	// to be stale in this case.
	ctx context.Context
	// Triangles has render data.
	Triangles []render.Triangle3
	// Sequence number of shape.
	Seq uint
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
