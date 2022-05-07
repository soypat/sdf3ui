package views

import (
	"github.com/soypat/three"
	"gonum.org/v1/gonum/spatial/r3"
)

func (v *shape3d) setCamera() {
	const maxDistance = 16
	size := bbSize(v.bb)
	sizeNorm := r3.Norm(size)
	center := bbCenter(v.bb)
	v.camera.SetFar(3 * sizeNorm * maxDistance)
	v.camera.SetNear(sizeNorm / 1e3)
	// ISO view looking at origin.
	camPos := r3.Add(center, r3.Vec{X: sizeNorm, Y: sizeNorm, Z: sizeNorm})
	v.camera.SetPosition(three.NewVector3(camPos.X, camPos.Y, camPos.Z))
	v.camera.LookAt(three.NewVector3(center.X, center.Y, center.Z))
	v.controls.SetTarget(three.NewVector3(center.X, center.Y, center.Z))
	v.controls.SetMaxDistance(sizeNorm * maxDistance)
	// Move mesh to certain location
	// mx, my, mz := v.shapeMesh.GetPosition().Coords()
	// v.shapeMesh.SetPosition(three.NewVector3(mx-center.X, my-center.Y, mz-center.Z))
}
