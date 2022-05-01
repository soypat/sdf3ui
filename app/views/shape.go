package views

import (
	"fmt"
	"math"
	"syscall/js"
	"vecty-templater-project/model"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/soypat/three"
	"github.com/soypat/three/vthree"
	"gonum.org/v1/gonum/spatial/r3"
)

type shape3d struct {
	vecty.Core

	shape       model.Shape3D
	shapeMesh   three.Mesh
	camera      three.PerspectiveCamera
	scene       three.Scene
	controls    three.TrackballControls
	near, far   float64
	renderedSeq int
}

func (v *shape3d) Render() vecty.ComponentOrHTML {
	b := &vthree.Basic{
		Init:    v.init,
		Animate: v.animate,
	}
	return elem.Div(
		b,
	)
}

func (v *shape3d) init(wgl three.WebGLRenderer) {
	elem := wgl.DomElement()
	js.Global().Set("webglel", elem)
	width := elem.Get("width").Float()
	height := elem.Get("height").Float()
	pixelRatio := js.Global().Get("devicePixelRatio").Float()
	wgl.SetPixelRatio(pixelRatio)
	wgl.SetSize(width, height, true)
	fmt.Println("webgl initialized with widthxheight", width, height)
	v.scene = three.NewScene()
	v.near = 0.2
	v.far = 100
	// Lights.
	dlight := three.NewDirectionalLight(three.NewColor("white"), 1) // lights, without lights everything will be dark!
	dlight.SetPosition(three.NewVector3(1000, 1000, 0))
	amblight := three.NewAmbientLight(three.NewColor("white"), 0.2)
	v.scene.Add(dlight)
	v.scene.Add(amblight)
	// v.scene.Add(three.NewFog(*three.NewColorHex(0x050505), v.far/2, v.far))
	// Camera.
	// ISO view looking at origin.
	v.camera = three.NewPerspectiveCamera(70, width/height, v.near, v.far)
	v.camera.SetPosition(three.NewVector3(v.far/2, v.far/2, v.far/2))
	v.camera.LookAt(three.NewVector3(0, 0, 0))

	// Controls.
	v.controls = three.NewTrackballControls(v.camera, wgl.DomElement())
	v.controls.SetMaxDistance(v.far * 1.4)

	v.renderShape(wgl)
	wgl.Render(v.scene, v.camera)
}

func (v *shape3d) animate(wgl three.WebGLRenderer) bool {
	v.renderShape(wgl)
	v.camera.SetFar(v.far)
	v.controls.Update()
	wgl.Render(v.scene, v.camera)
	return true
}

// SetShape sets the 3D shape.
func (v *shape3d) SetShape(shape model.Shape3D) {
	v.shape = shape
}

func (v *shape3d) renderShape(wgl three.WebGLRenderer) {
	if v.shape.Seq == uint(v.renderedSeq) {
		// Already rendered.
		fmt.Println("skipping render sequence", v.shape.Seq)
		return
	}
	if len(v.shape.Triangles) == 0 {
		fmt.Println("skipping render due to empty triangles")
		return
	}
	v.renderedSeq = int(v.shape.Seq)
	Nfaces := len(v.shape.Triangles)
	const faceLen = 3 * 3
	vertices := make([]float32, Nfaces*faceLen)
	normals := make([]float32, Nfaces*faceLen)
	var min, max r3.Vec
	for iface, face := range v.shape.Triangles {
		// vertices index of face.
		vertexStart := iface * faceLen
		n := face.Normal()
		for i := 0; i < 3; i++ {
			min = minElem(min, face.V[i])
			max = maxElem(max, face.V[i])
			vertexIdx := vertexStart + i*3
			vertices[vertexIdx] = float32(face.V[i].X)
			vertices[vertexIdx+1] = float32(face.V[i].Y)
			vertices[vertexIdx+2] = float32(face.V[i].Z)

			normals[vertexIdx] = float32(n.X)
			normals[vertexIdx+1] = float32(n.Y)
			normals[vertexIdx+2] = float32(n.Z)
		}
	}
	fmt.Println("finished getting ", Nfaces, " faces")
	geom := three.NewBufferGeometry()
	geom.SetAttribute("position", three.NewBufferAttribute(vertices, 3))
	geom.SetAttribute("normal", three.NewBufferAttribute(normals, 3))
	geom.ComputeBoundingSphere()
	material := three.NewMeshPhongMaterial(&three.MaterialParameters{
		Color:    three.NewColor("purple"),
		Specular: three.NewColor("gray"),
		Side:     three.DoubleSide,
	})
	mesh := three.NewMesh(geom, material)
	if v.shapeMesh.Truthy() {
		v.scene.Remove(v.shapeMesh)
		// v.shapeMesh.Call("destroy") // how to gracefully free memory?
	}
	size := r3.Sub(max, min)
	v.far = r3.Norm(size) * 8
	center := r3.Add(min, r3.Scale(0.5, size))
	v.camera.SetPosition(three.NewVector3(1.1*max.X, 1.1*max.Y, 1.1*max.Z))
	// v.camera.LookAt(three.NewVector3(center.X, center.Y, center.Z))
	v.controls.SetTarget(three.NewVector3(center.X, center.Y, center.Z))
	v.shapeMesh = mesh
	v.scene.Add(v.shapeMesh)
	fmt.Println("added new shape mesh")
}

// minElem return a vector with the minimum components of two vectors.
func minElem(a, b r3.Vec) r3.Vec {
	return r3.Vec{X: math.Min(a.X, b.X), Y: math.Min(a.Y, b.Y), Z: math.Min(a.Z, b.Z)}
}

// maxElem return a vector with the maximum components of two vectors.
func maxElem(a, b r3.Vec) r3.Vec {
	return r3.Vec{X: math.Max(a.X, b.X), Y: math.Max(a.Y, b.Y), Z: math.Max(a.Z, b.Z)}
}
