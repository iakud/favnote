package main

import (
	"github.com/iakud/coral"
)

func main() {
	loadBlogs()
	// parse template
	loadTemplate()
	// init router
	coral.Get("/favicon.ico", faviconIcoHandler)
	coral.Get("/static/(.*)", staticHandler)
	coral.Get("/img/(.*)", imageHandler)
	coral.Get("/article/(.*)", articleHandler)
	coral.Get("/", homeHandler)
	// run
	coral.Run("localhost:80")
}
