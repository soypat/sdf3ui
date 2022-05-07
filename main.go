package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"vecty-templater-project/model"

	"github.com/soypat/gwasm"
)

//go:generate go run ./assets/gen_assets.go

var (
	rendererServer = newRendererHandler()
	// Embed assets into binary.
	//go:embed assets/*
	assetsFS embed.FS
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("first and only argument must be file name of program")
	}
	filename := os.Args[1]
	if _, err := os.Stat(filename); err != nil {
		log.Fatal(err)
	}
	rendererServer.SetFileTarget(filename)
	go rendererServer.Start(context.Background())
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
	http.HandleFunc(model.ShapeEndpoint, rendererServer.server.serveShapeHTTP)
	log.Printf("listening on %v", model.HTTPBaseURL)
	log.Fatal(http.ListenAndServe(model.HTTPAddr, newBaseMiddleware(http.DefaultServeMux)))
}
