package store

import (
	"context"
	"fmt"
	"vecty-templater-project/app/store/actions"
	"vecty-templater-project/model"
)

var (
	Ctx            actions.Context
	shape          model.Shape3D
	shapeStaleFunc func()
	ServerReply    string
	Listeners      = newListenerRegistry()
)

func GetShape() model.Shape3D {
	return shape
}

func SetShape(s model.Shape3D) {
	if shapeStaleFunc != nil {
		shapeStaleFunc()
	}
	s.Ctx, shapeStaleFunc = context.WithCancel(context.Background())
	shape = s
}

func OnAction(action interface{}) {
	switch a := action.(type) {
	case *actions.GetShape:
		// fire shape updater listener
	case *actions.PageSelect:
		oldCtx := Ctx
		Ctx = actions.Context{
			Page:     a.Page,
			Referrer: &oldCtx,
		}

	case *actions.Back:
		Ctx = *Ctx.Referrer
	case *actions.Refresh:
		// do nothing, just fire listeners to refresh page.
	default:
		panic("unknown action selected!")
	}
	fmt.Printf("action %T", action)
	Listeners.Fire(action)
}
