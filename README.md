# sdf3ui
Browser based 3D model visualizer for use with soypat/sdf package. 3 minute Youtube demo and tutorial [here](https://www.youtube.com/watch?v=t-N9gOMKupw&ab_channel=PatWhitti).

![sdf3ui](https://user-images.githubusercontent.com/26156425/168592432-466dce19-6764-46a2-9218-3f49aed96dc2.gif)

#### Installation
Before installing make sure your `GOBIN` environment variable is set to a folder in your PATH.
```shell
go install github.com/soypat/sdf3ui@latest
```

### Example of usage
1. Create a `.go` extension file of any name in an empty folder. Copy the following contents into it. 
2. Run `go mod init my-awesome-module` in the directory followed by `go mod tidy` to download dependecies.
3. Run `sdf3ui thefile.go` in the directory.

```go
package main

import (
	"log"
	"os"

	"github.com/soypat/sdf"
	"github.com/soypat/sdf/form3"
	"github.com/soypat/sdf/render"
	"github.com/soypat/sdf3ui/uirender"
	"gonum.org/v1/gonum/spatial/r3"
)

func main() {
	var object sdf.SDF3
	object, _ = form3.Box(r3.Vec{4, 3, 2}, 1)
	cone, _ := form3.Cone(6, 1, 0, .1)
	union := sdf.Union3D(object, cone)
	union.SetMin(sdf.MinPoly(.4))
	object = union
	err := uirender.EncodeRenderer(os.Stdout, render.NewOctreeRenderer(object, 200))
	if err != nil {
		log.Fatal(err)
	}
}
```

![example_screenshot](https://user-images.githubusercontent.com/26156425/167276444-7898f12f-ff25-403a-a0fd-af5e48b2ad21.png)


### Build from source
If building from source run `go generate` to generate `assets` folder structure and contents.

To run app run `go run . example/mybolt.go` and head over to [http://localhost:8080](http://localhost:8080)

### Related projects
* [sdfx-ui](https://github.com/Yeicor/sdfx-ui) Real-time rendering using Faux-GL renderer for sdfx
