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
	imageFilePath := filepath.Join("blog/image", filePath)
	ctx.ServeFile(imageFilePath)
}

func blogHandler(ctx *coral.Context, name string) {
	if b, ok := allBlogs[name]; ok {
		renderTemplate(ctx, "blog", b)
	} else {
		ctx.NotFound(name)
	}
}

func homeHandler(ctx *coral.Context) {
	renderTemplate(ctx, "home", listBlogs)
}
