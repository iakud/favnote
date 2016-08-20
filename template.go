package main

import (
	"fmt"
	"github.com/iakud/coral"
	"html/template"
)

var tmplates map[string]*template.Template = make(map[string]*template.Template)

func loadTemplate() {
	tmplates["home"] = template.Must(template.ParseFiles("html/index.html", "html/layout.html"))
	tmplates["blog"] = template.Must(template.ParseFiles("html/blog.html", "html/layout.html"))
}

func renderTemplate(ctx *coral.Context, name string, data interface{}) error {
	if t, ok := tmplates[name]; ok {
		ctx.SetHeader("Content-Type", "text/html; charset=utf-8")
		return t.Execute(ctx.ResponseWriter, data)
	} else {
		return fmt.Errorf("template %q is undefined", name)
	}
}
