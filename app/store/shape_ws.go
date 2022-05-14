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

var (
	// prevent multiple requests from attempting to get shape at same time
	gettingShape = false
	wsConn       *websocket.Conn
)

func initShapeWS() error {
	if wsConn != nil {
		wsConn.Close(websocket.StatusAbnormalClosure, "client wanted to reinitialize")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost"+model.HTTPAddr+"/"+model.WSSubprotocol, &websocket.DialOptions{
		Subprotocols: []string{model.WSSubprotocol},
	})

	if err != nil {
		return err
	}
	c.SetReadLimit(model.WSReadLimit)
	wsConn = c
	return nil
}

func WebsocketShapeListen() {
	defer panic("WebSocketListen ended")
	for {
		err := initShapeWS()
		if err == nil {
			break
		}
		fmt.Println("websocket init failed. retry. ", err.Error())
	}
	fmt.Println("initialized websocket succesfully")
	var stat model.ServerStatus
	for {
		ctx, cancel := context.WithTimeout(context.Background(), model.WSTimeout)
		err := wsjson.Read(ctx, wsConn, &stat)
		cancel()
		if err != nil {
			fmt.Println("got websocket error", err)
			initShapeWS()
			continue
		}
		// if gettingShape {
		// 	continue
		// }
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
	if gettingShape {
		fmt.Println("cancel ForceUpdateShape. Shape already being fetched")
		return
	}
	gettingShape = true
	go func() {
		defer func() {
			gettingShape = false
		}()
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

		fmt.Println("Updated shape", shape.String())
		actions.Dispatch(&actions.Refresh{})
	}()
}
