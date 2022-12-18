package main

import (
	"log"
	"os"

	"github.com/soypat/sdf3ui/uirender"

	"github.com/soypat/sdf/form3/obj3/thread"
	"github.com/soypat/sdf/render"
)

func main() {
	const quality = 200
  
	b, _ := thread.Bolt(thread.BoltParms{
    Thread:      thread.ISO{D:16, P:2},
		Style:       thread.NutHex,
		Tolerance:   0.1,
		TotalLength: 80.,
		ShankLength: 0,
	}) 

	err := uirender.EncodeRenderer(os.Stdout, render.NewOctreeRenderer(b, quality))
	if err != nil {
		log.Fatal(err)
	}
}
