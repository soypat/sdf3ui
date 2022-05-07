package uirender

import (
	"bytes"
	"io"
	"reflect"
	"unsafe"

	"github.com/soypat/sdf/render"
)

func EncodeRenderer(dst io.Writer, src render.Renderer) error {
	t, err := render.RenderAll(src)
	if err != nil {
		return err
	}
	Nt := len(t)
	sizeTriangle := unsafe.Sizeof(render.Triangle3{})
	slice := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&t[0])),
		Len:  Nt * int(sizeTriangle),
		Cap:  Nt * int(sizeTriangle),
	}
	byteSlice := *(*[]byte)(unsafe.Pointer(&slice))
	_, err = io.Copy(dst, bytes.NewReader(byteSlice))
	return err
}

type Decoder struct {
	t []render.Triangle3
	n int64
}

func DecodeAll(src io.Reader) ([]render.Triangle3, error) {
	d := DecodeRenderer(src).(*Decoder)
	return d.t, nil
}
func DecodeRenderer(src io.Reader) render.Renderer {
	b, err := io.ReadAll(src)
	if err != nil {
		panic(err)
	}
	sizeTriangle := int(unsafe.Sizeof(render.Triangle3{}))
	nt := len(b) / sizeTriangle
	if len(b) != nt*sizeTriangle {
		panic("bad result length Decoder")
	}
	slice := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&b[0])),
		Len:  nt,
		Cap:  nt,
	}
	triSlice := *(*[]render.Triangle3)(unsafe.Pointer(&slice))
	return &Decoder{t: triSlice}
}

func (d *Decoder) ReadTriangles(dst []render.Triangle3) (nt int, err error) {
	if d.n == int64(len(d.t)) {
		return 0, io.EOF
	}
	nt = copy(dst, d.t[d.n:])
	d.n += int64(nt)
	return nt, nil
}
