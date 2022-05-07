package store

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/soypat/sdf3ui/app/store/actions"
	"github.com/soypat/sdf3ui/model"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var wsConn *websocket.Conn

func initShapeWS() error {
	fmt.Println("init socket")
	if wsConn != nil {
		wsConn.Close(websocket.StatusAbnormalClosure, "client wanted to reinitialize")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost"+model.HTTPAddr+"/"+model.WSSubprotocol, &websocket.DialOptions{
		Subprotocols: []string{model.WSSubprotocol},
	})

	if err != nil {
		fmt.Println("websocket initialization failed:", err.Error())
		return err
	}
	c.SetReadLimit(model.MaxRenderSize)
	wsConn = c
	fmt.Println("initialized websocket succesfully")
	return nil
}

func WebsocketShapeListen() {
	defer panic("WebSocketListen ended")
	for {
		err := initShapeWS()
		if err == nil {
			break
		}
		fmt.Println("websocket init failed. retry")
	}

	var stat model.ServerStatus
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Second)
		err := wsjson.Read(ctx, wsConn, &stat)
		cancel()
		if err != nil {
			fmt.Println("got websocket error", err)
			initShapeWS()
			continue
		}
		currentShape := GetShape()
		if currentShape.Seq >= stat.ShapeSeq {
			// shape is not stale, do nothing
			continue
		}
		// Shape is stale- update current shape.
		ForceUpdateShape()
	}
}

func ForceUpdateShape() {
	go func() {
		defer TimeIt("ForceUpdateShape")()
		resp, err := http.Get(model.ShapeEndpoint)
		if err != nil {
			fmt.Println("error during GET updating shape", err)
			return
		}
		var gotShape model.Shape3D

		err = gotShape.Decode(resp.Body)
		if err != nil {
			fmt.Println("error during shape decoding:", err)
			return
		}
		SetShape(gotShape)
		fmt.Println("shape update success. amount of triangles streamed:", len(shape.Triangles))
		actions.Dispatch(&actions.Refresh{})
	}()
}
