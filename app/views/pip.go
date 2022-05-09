package views

import (
	"errors"
	"syscall/js"

	"github.com/hexops/vecty"
)

// PIPWindow implements https://developer.mozilla.org/en-US/docs/Web/API/Picture-in-Picture_API
type PIPWindow struct {
	// Only Elem supported is HTMLVideoElement
	Elem     js.Value                                    `vecty:"prop"`
	OnResize func(width, height float64, e *vecty.Event) `vecty:"prop"`
}

// RequestPIP is called when the PIP
// functionality is desired.
func (p *PIPWindow) RequestPIP() error {
	promise := p.Elem.Call("requestPictureInPicture")
	result := promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		pipWindow := args[0]
		pipWindow.Set("onresize", js.FuncOf(p.onresize))
		return nil
	}))
	var err error
	result.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		msg := args[0].String()
		if msg != "" {
			err = errors.New(msg)
		}
		return nil
	}))
	return err
}

func (p *PIPWindow) onresize(this js.Value, args []js.Value) interface{} {
	ev := args[0]
	pipWindow := ev.Get("target")
	width := pipWindow.Get("width").Float()
	height := pipWindow.Get("height").Float()
	p.OnResize(width, height, &vecty.Event{Target: pipWindow, Value: ev})
	return nil
}
