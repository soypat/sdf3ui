package main

import (
	"fmt"
	"time"

	"github.com/soypat/sdf3ui/app/store"
	"github.com/soypat/sdf3ui/app/store/actions"

	"github.com/soypat/sdf3ui/app/views"

	"github.com/hexops/vecty"
	"github.com/soypat/gwasm"
	"github.com/soypat/three"
)

func main() {
	Message := "Welcome!"
	gwasm.AddScript("https://threejs.org/build/three.js", "THREE", 3*time.Second)

	fmt.Println("if assets/js/trackball_controls.js fails to get please run `go generate` in sdf3ui base directory to generate assets")
	gwasm.AddScript("assets/js/trackball_controls.js", "TrackballControls", time.Second)
	err := three.Init()
	if err != nil {
		Message = "three.js not found!"
	}
	go store.WebsocketShapeListen()
	// OnAction must be registered before any storage manipulation.
	actions.Register(store.OnAction)
	// Get Latest shape.
	store.ForceUpdateShape()
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
