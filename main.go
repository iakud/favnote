package main

import (
	"github.com/iakud/favor/web"
)

func main() {
	loadBlogs()
	// parse template
	loadTemplate()
	// init router
	web.Get("/favicon.ico", faviconIcoHandler)
	web.Get("/static/(.*)", staticHandler)
	web.Get("/img/(.*)", imageHandler)
	web.Get("/article/(.*)", articleHandler)
	web.Get("/", homeHandler)
	// run
	web.Run("localhost:80")
}
