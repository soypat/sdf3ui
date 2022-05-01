package views

import (
	"vecty-templater-project/app/store"
	"vecty-templater-project/app/store/actions"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
)

type Body struct {
	vecty.Core
	Ctx  actions.Context `vecty:"prop"`
	Info string          `vecty:"prop"`
}

func (b *Body) Render() vecty.ComponentOrHTML {
	var mainContent vecty.MarkupOrChild
	switch b.Ctx.Page {
	case actions.PageLanding:
		mainContent = elem.Div(
			elem.Strong(vecty.Text(b.Info)),
			elem.Div(elem.Button(
				vecty.Markup(event.Click(b.newItem)),
				vecty.Text("Refresh Shape3D"),
			)),
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
