package model

import (
	"context"

	"github.com/soypat/sdf/render"
)

const (
	HTTPServerAddr = ":8080"
	// Websocket sub protocol.
	WSSubprotocol = "sdf3ui"
	megabyte      = 1000 * 1000
	MaxRenderSize = 30 * megabyte
)

// Shape3D contains 3D shape information.
type Shape3D struct {
	// Ctx stores shape status. This context is cancelled when
	// a newer shape is available. The previous shape is said
	// to be stale in this case.
	Ctx context.Context `json:"-"`
	// Triangles has render data.
	Triangles []render.Triangle3
	// Sequence number of shape.
	Seq uint
}
