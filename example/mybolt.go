package main

import (
	"log"
	"os"
	"vecty-templater-project/uirender"

	"github.com/soypat/sdf/form3/obj3"
	"github.com/soypat/sdf/render"
)

func main() {
	b, _ := obj3.Bolt(obj3.BoltParms{
		Thread:      "npt_1/2",
		Style:       obj3.CylinderHex,
		TotalLength: 10,
		ShankLength: 3,
	})
	err := uirender.EncodeRenderer(os.Stdout, render.NewOctreeRenderer(b, 200))
	if err != nil {
		log.Fatal(err)
	}
}
