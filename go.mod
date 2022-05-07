module github.com/soypat/sdf3ui

go 1.17

require (
	github.com/fsnotify/fsnotify v1.5.4
	github.com/hexops/vecty v0.6.0
	github.com/soypat/gwasm v0.0.10
	github.com/soypat/sdf v0.0.0-20220507034430-c26521433a8b
	github.com/soypat/three v0.0.0-20220501183759-e2d2208ad9b5
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	gonum.org/v1/gonum v0.11.0
	nhooyr.io/websocket v1.8.7
)

require (
	github.com/klauspost/compress v1.10.3 // indirect
	github.com/soypat/rebed v0.2.3 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
)

replace github.com/soypat/gwasm => ../gwasm
