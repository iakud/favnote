package main

import (
	"github.com/iakud/favor/web"
	"path/filepath"
)

func faviconIcoHandler(ctx *web.Context) {
	ctx.ServeFile("static/favicon.ico")
}

func staticHandler(ctx *web.Context, filePath string) {
	staticFilePath := filepath.Join("static", filePath)
	ctx.ServeFile(staticFilePath)
}

func imageHandler(ctx *web.Context, filePath string) {
	imageFilePath := filepath.Join("static", filePath)
	ctx.ServeFile(imageFilePath)
}

func articleHandler(ctx *web.Context, name string) {
	if b, ok := articles[name]; ok {
		renderTemplate(ctx, "article", b)
	} else {
		ctx.NotFound(name)
	}
}

func homeHandler(ctx *web.Context) {
	renderTemplate(ctx, "home", sortedArticles)
}
