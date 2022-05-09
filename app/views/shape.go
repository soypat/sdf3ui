package views

import (
	"fmt"
	"math"
	"syscall/js"
	"time"

	"github.com/soypat/sdf3ui/app/store"
	"github.com/soypat/sdf3ui/model"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/soypat/sdf/render"
	"github.com/soypat/three"
	"github.com/soypat/three/vthree"
	"gonum.org/v1/gonum/spatial/r3"
)

type shape3d struct {
	vecty.Core

	shape model.Shape3D
	// Bounding box of shape
	bb            r3.Box
	shapeMesh     three.Mesh
	camera        three.PerspectiveCamera
	scene         three.Scene
	controls      three.TrackballControls
	pip           PIPWindow
	width, height float64
	lastResize    time.Time
	renderedSeq   int
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
	pixelRatio := js.Global().Get("devicePixelRatio").Float()
	wgl.SetPixelRatio(pixelRatio)
	v.pip.Elem = js.Global().Get("document") // wgl.DomElement()
	v.pip.OnResize = func(width, height float64, e *vecty.Event) {
		v.width = width
		v.height = height
	}
	v.setSize(wgl)
	v.scene = three.NewScene()
	// Lights.  without lights everything will be dark!
	dlight := three.NewDirectionalLight(three.NewColor("white"), 1)
	dlight.SetPosition(three.NewVector3(1000, 1000, 0))
	amblight := three.NewAmbientLight(three.NewColor("white"), 0.2)
	dlight2 := three.NewDirectionalLight(three.NewColor("lightblue"), .5)
	dlight2.SetPosition(three.NewVector3(0, -1000, 1000))
	v.scene.Add(dlight)
	v.scene.Add(amblight)
	v.scene.Add(dlight2)

	// Camera.
	// ISO view looking at origin.
	v.camera = three.NewPerspectiveCamera(70, 4/3, 0.1, 2000)

	// Controls.
	v.controls = three.NewTrackballControls(v.camera, wgl.DomElement())

	v.renderShape(wgl)
	wgl.Render(v.scene, v.camera)
}

func (v *shape3d) animate(wgl three.WebGLRenderer) bool {
	v.setSize(wgl)
	v.renderShape(wgl)
	v.controls.Update()
	wgl.Render(v.scene, v.camera)
	return true
}

func (v *shape3d) setSize(wgl three.WebGLRenderer) {
	if time.Since(v.lastResize) < time.Second {
		return
	}
	elem := wgl.DomElement()
	currentWidth := elem.Get("width").Float()
	currentHeight := elem.Get("height").Float()
	if currentWidth != v.width || currentHeight != v.height {
		if v.width == 0 && v.height == 0 {
			v.width = currentWidth
			v.height = currentHeight
		}
		v.lastResize = time.Now()
		wgl.SetSize(v.width, v.height, true)
		fmt.Println("webgl initialized with widthxheight", v.width, v.height)
	}
}

// SetShape sets the 3D shape.
func (v *shape3d) SetShape(shape model.Shape3D) {
	v.shape = shape
}

func (v *shape3d) renderShape(wgl three.WebGLRenderer) {
	if v.shape.Seq == uint(v.renderedSeq) {
		// Already rendered.
		// fmt.Println("skipping render sequence", v.shape.Seq)
		return
	}
	if len(v.shape.Triangles) == 0 {
		// fmt.Println("skipping render due to empty triangles")
		return
	}
	defer v.setCamera()
	v.renderedSeq = int(v.shape.Seq)
	mesh, box := makeShapeMesh(v.shape.Triangles)
	v.bb = box
	// points, _ := makePointMesh(v.shape.Triangles) // for debugging
	// mesh.Add(points)
	if v.shapeMesh.Truthy() {
		v.scene.Remove(v.shapeMesh)
		// v.shapeMesh.Call("dispose") // does this free all memory?
	}

	v.shapeMesh = mesh
	v.scene.Add(v.shapeMesh)
	fmt.Println("added new shape mesh")
}

func (v *shape3d) setCamera() {
	defer store.TimeIt("shape3d.setCamera")()
	size := bbSize(v.bb)
	sizeNorm := r3.Norm(size)
	center := bbCenter(v.bb)
	far := 4 * sizeNorm
	v.camera.SetFar(far)
	v.camera.SetNear(sizeNorm / 1e3)
	// ISO view looking at origin.
	camPos := r3.Add(center, r3.Vec{X: sizeNorm, Y: sizeNorm, Z: sizeNorm})
	v.camera.SetPosition(three.NewVector3(camPos.X, camPos.Y, camPos.Z))
	v.camera.LookAt(three.NewVector3(center.X, center.Y, center.Z))
	v.controls.SetTarget(three.NewVector3(center.X, center.Y, center.Z))
	v.controls.SetMaxDistance(2 * sizeNorm)

	v.camera.UpdateProjectionMatrix()
}

// minElem return a vector with the minimum components of two vectors.
func minElem(a, b r3.Vec) r3.Vec {
	return r3.Vec{X: math.Min(a.X, b.X), Y: math.Min(a.Y, b.Y), Z: math.Min(a.Z, b.Z)}
}

// maxElem return a vector with the maximum components of two vectors.
func maxElem(a, b r3.Vec) r3.Vec {
	return r3.Vec{X: math.Max(a.X, b.X), Y: math.Max(a.Y, b.Y), Z: math.Max(a.Z, b.Z)}
}

func bbSize(bb r3.Box) r3.Vec {
	return r3.Sub(bb.Max, bb.Min)
}

func bbCenter(bb r3.Box) r3.Vec {
	return r3.Add(bb.Min, r3.Scale(0.5, bbSize(bb)))
}

func makeShapeMesh(t []render.Triangle3) (three.Mesh, r3.Box) {
	defer store.TimeIt("shape3D.makeShapeMesh")()
	Nfaces := len(t)
	const faceLen = 3 * 3
	vertices := make([]float32, Nfaces*faceLen)
	normals := make([]float32, Nfaces*faceLen)
	var min, max r3.Vec
	for iface, face := range t {
		// vertices index of face.
		vertexStart := iface * faceLen
		n := face.Normal()
		nx := float32(n.X)
		ny := float32(n.Y)
		nz := float32(n.Z)
		for i := 0; i < 3; i++ {
			min = minElem(min, face.V[i])
			max = maxElem(max, face.V[i])
			vertexIdx := vertexStart + i*3
			vertices[vertexIdx] = float32(face.V[i].X)
			vertices[vertexIdx+1] = float32(face.V[i].Y)
			vertices[vertexIdx+2] = float32(face.V[i].Z)

			normals[vertexIdx] = nx
			normals[vertexIdx+1] = ny
			normals[vertexIdx+2] = nz
		}
	}
	fmt.Println("finished getting ", Nfaces, " faces")
	geom := three.NewBufferGeometry()
	geom.SetAttribute("position", three.NewBufferAttribute(vertices, 3))
	geom.SetAttribute("normal", three.NewBufferAttribute(normals, 3))
	geom.ComputeBoundingSphere()
	material := three.NewMeshPhongMaterial(&three.MaterialParameters{
		Color:    three.NewColor("chocolate"),
		Specular: three.NewColor("gray"),
		Side:     three.FrontSide,
	})
	mesh := three.NewMesh(geom, material)
	return mesh, r3.Box{Min: min, Max: max}
}

func makePointMesh(t []render.Triangle3) (three.Points, r3.Box) {
	Nfaces := len(t)
	const faceLen = 3 * 3
	vertices := make([]float32, Nfaces*faceLen)
	var min, max r3.Vec
	for iface, face := range t {
		// vertices index of face.
		vertexStart := iface * faceLen
		vertices[vertexStart+0] = float32(face.V[0].X)
		vertices[vertexStart+1] = float32(face.V[0].Y)
		vertices[vertexStart+2] = float32(face.V[0].Z)

		vertices[vertexStart+3] = float32(face.V[1].X)
		vertices[vertexStart+4] = float32(face.V[1].Y)
		vertices[vertexStart+5] = float32(face.V[1].Z)

		vertices[vertexStart+6] = float32(face.V[2].X)
		vertices[vertexStart+7] = float32(face.V[2].Y)
		vertices[vertexStart+8] = float32(face.V[2].Z)
	}
	fmt.Println("finished getting ", Nfaces, " faces")
	geom := three.NewBufferGeometry()
	geom.SetAttribute("position", three.NewBufferAttribute(vertices, 3))

	geom.ComputeBoundingSphere()
	mesh := three.NewPoints(geom, three.NewPointsMaterial(three.MaterialParameters{
		Color: three.NewColor("red"),
		Size:  .1,
	}))
	return mesh, r3.Box{Min: min, Max: max}
}

func (v *shape3d) requestPIP() error {
	return v.pip.RequestPIP()
}
