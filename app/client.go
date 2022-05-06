package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"vecty-templater-project/app/store"
	"vecty-templater-project/app/store/actions"
	"vecty-templater-project/model"

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
	updateShape()
	go updateShape()
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
		updateShape()
	})
}

func updateShape() {
	resp, err := http.Get(model.ShapeEndpoint)
	if err != nil {
		fmt.Println("getting new shape:", err)
		return
	}
	var shape model.Shape3D
	err = json.NewDecoder(resp.Body).Decode(shape)
	if err != nil {
		fmt.Println("decoding new shape:", err)
		return
	}

	fmt.Println("updateShape success, got shape of size ", len(shape.Triangles))
	store.SetShape(shape)
	store.ServerReply = fmt.Sprintf("Shape sequence %d", shape.Seq)
	actions.Dispatch(&actions.Refresh{})
}

// var wsConn *websocket.Conn

// func initShapeWS() {
// 	fmt.Println("init socket")
// 	if wsConn != nil {
// 		wsConn.Close(websocket.StatusAbnormalClosure, "client wanted to reinitialize")
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
// 	defer cancel()

// 	c, _, err := websocket.Dial(ctx, "ws://localhost"+model.HTTPAddr+"/"+model.WSSubprotocol, &websocket.DialOptions{
// 		Subprotocols: []string{model.WSSubprotocol},
// 	})

// 	if err != nil {
// 		fmt.Println("websocket initialization failed:", err.Error())
// 		return
// 	}
// 	c.SetReadLimit(model.MaxRenderSize)
// 	wsConn = c
// 	fmt.Println("initialized websocket succesfully")
// }
