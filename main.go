package main

import (
	"github.com/iakud/coral"
)

func main() {
	loadAndWatchingBlogs()
	// parse template
	loadTemplate()
	// init router
	coral.Get("/favicon.ico", faviconIcoHandler)
	coral.Get("/static/(.*)", staticHandler)
	coral.Get("/img/(.*)", imageHandler)
	coral.Get("/blog/(.*)", blogHandler)
	coral.Get("/", homeHandler)
	// run
	coral.Run("localhost:80")
}
