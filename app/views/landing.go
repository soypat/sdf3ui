package views

import (
	"vecty-templater-project/model"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

var canvas = &shape3d{}

type Landing struct {
	vecty.Core
	Shape model.Shape3D `vecty:"prop"`
}

func (l *Landing) Render() vecty.ComponentOrHTML {
	canvas.SetShape(l.Shape)
	return elem.Div(
		canvas,
	)
}
