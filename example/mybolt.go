package main

import (
	"log"
	"os"
	"vecty-templater-project/uirender"

	"github.com/soypat/sdf/form3/obj3"
	"github.com/soypat/sdf/render"
)

func main() {
	const quality = 200
	b, _ := obj3.Bolt(obj3.BoltParms{
		Thread:      "M16x2",
		Style:       obj3.CylinderHex,
		Tolerance:   0.1,
		TotalLength: 60.0,
		ShankLength: 10.0,
	})
	err := uirender.EncodeRenderer(os.Stdout, render.NewOctreeRenderer(b, quality))
	if err != nil {
		log.Fatal(err)
	}
}
