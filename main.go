package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"vecty-templater-project/model"

	"github.com/soypat/gwasm"
	"github.com/soypat/sdf/form3/obj3"
	"github.com/soypat/sdf/render"
)

//go:generate go run ./assets/gen_assets.go

var (
	rendererServer = newRendererHandler()
	// Embed assets into binary.
	//go:embed assets/*
	assetsFS embed.FS
)

func main() {
	b, _ := obj3.Bolt(obj3.BoltParms{
		Thread:      "npt_1/2",
		Style:       obj3.CylinderHex,
		TotalLength: 30,
		ShankLength: 5,
	})
	t, err := render.RenderAll(render.NewOctreeRenderer(b, 200))
	if err != nil {
		log.Fatal(err)
	}
	rendererServer.server.SetShape(t)
	// Register http server for serving WASM and other resources.
	wsm, err := gwasm.NewWASMHandler("app", nil)
	if err != nil {
		log.Fatal("NewWASMHandler failed", err)
	}
	wsm.WASMReload = true
	wsm.SetOutput(os.Stdout)
	http.Handle("/", wsm)
	http.Handle("/"+model.WSSubprotocol, &rendererServer.server)
	http.Handle("/assets/", http.FileServer(http.FS(assetsFS)))
	log.Printf("listening on http://[::]%v", model.HTTPServerAddr)
	log.Fatal(http.ListenAndServe(model.HTTPServerAddr, nil))
}
