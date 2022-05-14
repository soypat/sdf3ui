package views

import (
	"fmt"

	"github.com/soypat/sdf3ui/app/store"
	"github.com/soypat/sdf3ui/app/store/actions"
	"github.com/soypat/three"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type Body struct {
	vecty.Core
	Ctx  actions.Context `vecty:"prop"`
	Info string          `vecty:"prop"`

	stlNameInput *vecty.HTML
}

func (b *Body) Render() vecty.ComponentOrHTML {
	var mainContent vecty.MarkupOrChild
	switch b.Ctx.Page {
	case actions.PageLanding:
		b.stlNameInput = elem.Input(
			vecty.Markup(vecty.Attribute("placeholder", "stl filename")),
		)
		mainContent = elem.Div(
			elem.Strong(vecty.Text(b.Info)),
			elem.Div(
				elem.Button(
					vecty.Markup(event.Click(b.centerShape)),
					vecty.Text("Center Shape"),
				),
				elem.Button(
					vecty.Markup(event.Click(b.newItem)),
					vecty.Text("Refresh Shape3D"),
				),
				b.stlNameInput,
				elem.Button(
					vecty.Markup(event.Click(b.downloadSTL)),
					vecty.Text("Save STL in working directory"),
				),

				// elem.Button(
				// 	vecty.Markup(event.Click(b.requestPIP)),
				// 	vecty.Text("PIP (unsupported)"),
				// ),
			),
			&Landing{
				Shape: store.GetShape(),
			},
		)
	default:
		panic("unknown Page")
	}
	return elem.Body(
		vecty.If(b.Ctx.Referrer != nil, elem.Div(
			elem.Button(
				vecty.Markup(event.Click(b.backButton)),
				vecty.Text("Back"),
			))),
		mainContent,
	)
}

func (b *Body) backButton(*vecty.Event) {
	actions.Dispatch(&actions.Back{})
}

func (b *Body) newItem(*vecty.Event) {
	actions.Dispatch(&actions.GetShape{})
}

func (b *Body) downloadSTL(*vecty.Event) {
	filename := b.stlNameInput.Node().Get("value").String()
	if filename == "" {
		filename = "sdf3ui_output.stl"
	}
	store.SaveRemoteSTL(filename)
	actions.Dispatch(&actions.DownloadShapeSTL{Shape: store.GetShape()})
}

func (b *Body) centerShape(*vecty.Event) {
	v := bbCenter(canvas.bb)
	canvas.controls.SetTarget(three.NewVector3(v.X, v.Y, v.Z))
}

func (b *Body) requestPIP(*vecty.Event) {
	err := canvas.requestPIP()
	if err != nil {
		fmt.Println("requesting PIP:", err)
	}
}
