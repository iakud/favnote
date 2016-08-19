package main

import (
	"github.com/iakud/coral"
	"path/filepath"
)

func faviconIcoHandler(ctx *coral.Context) {
	ctx.ServeFile("static/favicon.ico")
}

func staticHandler(ctx *coral.Context, filePath string) {
	staticFilePath := filepath.Join("static", filePath)
	ctx.ServeFile(staticFilePath)
}

func imageHandler(ctx *coral.Context, filePath string) {
	imageFilePath := filepath.Join("static", filePath)
	ctx.ServeFile(imageFilePath)
}

func articleHandler(ctx *coral.Context, name string) {
	if b, ok := articles[name]; ok {
		renderTemplate(ctx, "article", b)
	} else {
		ctx.NotFound(name)
	}
}

func homeHandler(ctx *coral.Context) {
	renderTemplate(ctx, "home", sortedArticles)
}
