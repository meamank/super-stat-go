package superstatgo

import (
	"embed"
	"io/fs"
	"net/http"
)

var webAssets embed.FS

func GetFileSystem() http.FileSystem {
	subFS, err := fs.Sub(webAssets, "web/dist")

	if err != nil {
		panic(err)
	}

	return http.FS(subFS)
}
