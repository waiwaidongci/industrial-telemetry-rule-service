package httpadapter

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/*
var consoleAssets embed.FS

func (s *Server) consoleHandler() http.Handler {
	assets, err := fs.Sub(consoleAssets, "web")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(assets))
}
