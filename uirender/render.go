package uirender

import (
	"encoding/gob"
	"io"

	"github.com/soypat/sdf/render"
)

func EncodeRenderer(dst io.Writer, src render.Renderer) error {
	t, err := render.RenderAll(src)
	if err != nil {
		return err
	}
	return gob.NewEncoder(dst).Encode(&t)
}

type Decoder struct {
	t []render.Triangle3
	n int64
}

func DecodeAll(src io.Reader) ([]render.Triangle3, error) {
	d, err := DecodeRenderer(src)
	if err != nil {
		return nil, err
	}

	return d.(*Decoder).t, nil
}
func DecodeRenderer(src io.Reader) (render.Renderer, error) {
	var t []render.Triangle3
	err := gob.NewDecoder(src).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &Decoder{t: t}, nil
}

func (d *Decoder) ReadTriangles(dst []render.Triangle3) (nt int, err error) {
	if d.n == int64(len(d.t)) {
		return 0, io.EOF
	}
	nt = copy(dst, d.t[d.n:])
	d.n += int64(nt)
	return nt, nil
}
