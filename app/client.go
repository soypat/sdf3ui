package main

import (
	"fmt"
	"time"

	"vecty-templater-project/app/store"
	"vecty-templater-project/app/store/actions"

	"vecty-templater-project/app/views"

	"github.com/hexops/vecty"
	"github.com/soypat/gwasm"
	"github.com/soypat/three"
)

func main() {
	Message := "Welcome!"
	gwasm.AddScript("https://threejs.org/build/three.js", "THREE", time.Second)

	fmt.Println("if assets/js/trackball_controls.js fails to get please run `go generate` in sdf3ui base directory to generate assets")
	gwasm.AddScript("assets/js/trackball_controls.js", "TrackballControls", time.Second)
	err := three.Init()
	if err != nil {
		Message = "three.js not found!"
	}
	go store.WebsocketShapeListen()
	// OnAction must be registered before any storage manipulation.
	actions.Register(store.OnAction)

	addShapeListener()

	body := &views.Body{
		Ctx:  store.Ctx,
		Info: Message,
	}
	store.Listeners.Add(body, func(interface{}) {
		body.Ctx = store.Ctx
		body.Info = store.ServerReply
		vecty.Rerender(body)
	})
	vecty.RenderBody(body)
}

// addShapeListener
func addShapeListener() {
	const key = "shape3"
	// defer initShapeWS()
	store.Listeners.Add(nil, func(action interface{}) {
		if _, ok := action.(*actions.GetShape); !ok {
			// Only render new shape on RefreshShape action.
			return
		}
		store.ForceUpdateShape()
	})
}
