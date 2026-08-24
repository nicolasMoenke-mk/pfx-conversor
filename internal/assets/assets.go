package assets

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var staticEmbed embed.FS

func GetCoreFile() http.FileSystem {

	fsys, err := fs.Sub(staticEmbed, "static")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}
