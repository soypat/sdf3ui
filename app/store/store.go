package store

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"vecty-templater-project/app/store/actions"
	"vecty-templater-project/model"

	"github.com/soypat/gwasm"
	"github.com/soypat/sdf/render"
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
	var ctx context.Context
	ctx, shapeStaleFunc = context.WithCancel(context.Background())
	s.SetContext(ctx)
	shape = s
}

func OnAction(action interface{}) {
	switch a := action.(type) {
	case *actions.GetShape:
		ForceUpdateShape()
		// fire shape updater listener
	case *actions.PageSelect:
		oldCtx := Ctx
		Ctx = actions.Context{
			Page:     a.Page,
			Referrer: &oldCtx,
		}

	case *actions.DownloadShapeSTL:
		r, w := io.Pipe()
		defer r.Close()
		go func() {
			defer w.Close()
			err := render.WriteSTL(w, a.Shape.Triangles)
			if err != nil {
				fmt.Println("writing STL to stream:", err)
			}
		}()
		err := gwasm.DownloadStream("shape.stl", "", r)
		if err != nil {
			fmt.Println("downloading STL:", err)
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

func SaveRemoteSTL(filename string) {
	go func() {
		resp, err := http.Get(model.HTTPBaseURL + model.SaveSTLEndpoint + "?name=" + filename)
		if err != nil || resp.StatusCode != 200 {
			msg, _ := io.ReadAll(resp.Body)
			fmt.Println("Saving remote STL error ", resp.StatusCode, string(msg), err)
		}
	}()
}
