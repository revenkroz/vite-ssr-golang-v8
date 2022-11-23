package main

import (
	"embed"
	"github.com/revenkroz/vite-ssr-golang/pkg"
	"io/fs"
)

//go:embed all:dist/client
var frontendDist embed.FS

//go:embed all:dist/server
var serverDist embed.FS

func main() {
	fsysFrontend, _ := fs.Sub(frontendDist, "dist/client")
	fsysServer, _ := fs.Sub(serverDist, "dist/server")

	pkg.RunBlocking(pkg.FrontendBuild{
		FrontendDist: fsysFrontend,
		ServerDist:   fsysServer,
	})
}
