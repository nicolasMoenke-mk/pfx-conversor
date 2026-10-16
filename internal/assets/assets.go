package assets

import (
	"embed"
	"io/fs"
	"net/http"
)

// staticEmbed embute a pasta static (index.html e demais arquivos) dentro
// do próprio binário em tempo de compilação é o que torna o programa um
// único .exe autossuficiente, sem precisar desses arquivos soltos no disco
// de quem for rodar.

//go:embed static
var staticEmbed embed.FS

// GetCoreFile devolve o conteúdo embutido pronto pra ser servido via
// http.FileServer. fs.Sub reancora a raiz em "static", pra servir os
// arquivos direto em "/" (ex.: "/index.html") em vez de "/static/index.html".
// O erro só ocorreria se "static" não existisse dentro do embed — algo que
// só pode acontecer por engano em tempo de compilação, nunca em produção,
// por isso o panic em vez de propagar o erro.
func GetCoreFile() http.FileSystem {

	fsys, err := fs.Sub(staticEmbed, "static")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}
