package model

import (
	"io"
	"testing"

	"github.com/soypat/sdf/render"
)

func TestShapeCodec(t *testing.T) {
	var s1 Shape3D
	s1.Seq = 3
	s1.Triangles = make([]render.Triangle3, 1)
	r, w := io.Pipe()
	go func() {
		err := s1.Encode(w)
		if err != nil {
			t.Error("encoding", err)
		}
	}()
	var dst Shape3D
	err := dst.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	if s1.Seq != dst.Seq {
		t.Error("Sequence number not match after encoding")
	}
	if s1.Triangles[0] != dst.Triangles[0] {
		t.Error("triangles not match after encoding")
	}
}
