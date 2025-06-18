module github.com/opensraph/sraph/surface

go 1.24.3

require (
	github.com/opensraph/sraph/gpu v0.0.0
	github.com/rajveermalviya/go-wayland/wayland v0.0.0-20230130181619-0ad78d1310b2
)

require golang.org/x/sys v0.33.0 // indirect

replace github.com/opensraph/sraph/gpu => ../gpu
