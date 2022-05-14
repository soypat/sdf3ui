package main

import (
	"log"
	"os"

	"github.com/soypat/sdf3ui/uirender"

	"github.com/soypat/sdf/form3/obj3"
	"github.com/soypat/sdf/render"
)

func main() {
	const quality = 200
	b, _ := obj3.Bolt(obj3.BoltParms{
		Thread:      "M16x2",
		Style:       obj3.CylinderHex,
		Tolerance:   0.1,
		TotalLength: 80.,
		ShankLength: 0.0,
	})
	err := uirender.EncodeRenderer(os.Stdout, render.NewOctreeRenderer(b, quality))
	if err != nil {
		log.Fatal(err)
	}
}
