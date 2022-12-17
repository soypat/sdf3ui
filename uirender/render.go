package uirender

import (
	"encoding/gob"
	"io"

	"github.com/soypat/sdf/render"
	"gonum.org/v1/gonum/spatial/r3"
)

func EncodeRenderer(dst io.Writer, src render.Renderer) error {
	t, err := render.RenderAll(src)
	if err != nil {
		return err
	}
	return gob.NewEncoder(dst).Encode(&t)
}

type Decoder struct {
	t []r3.Triangle
	n int64
}

func DecodeAll(src io.Reader) ([]r3.Triangle, error) {
	d, err := DecodeRenderer(src)
	if err != nil {
		return nil, err
	}

	return d.(*Decoder).t, nil
}
func DecodeRenderer(src io.Reader) (render.Renderer, error) {
	var t []r3.Triangle
	err := gob.NewDecoder(src).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &Decoder{t: t}, nil
}

func (d *Decoder) ReadTriangles(dst []r3.Triangle) (nt int, err error) {
	if d.n == int64(len(d.t)) {
		return 0, io.EOF
	}
	nt = copy(dst, d.t[d.n:])
	d.n += int64(nt)
	return nt, nil
}
