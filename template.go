package main

import (
	"fmt"
	"github.com/iakud/favor/web"
	"html/template"
)

var tmplates map[string]*template.Template = make(map[string]*template.Template)

func loadTemplate() {
	tmplates["home"] = template.Must(template.ParseFiles("html/index.html", "html/layout.html"))
	tmplates["article"] = template.Must(template.ParseFiles("html/article.html", "html/layout.html"))
}

func renderTemplate(ctx *web.Context, name string, data interface{}) error {
	if t, ok := tmplates[name]; ok {
		ctx.SetHeader("Content-Type", "text/html; charset=utf-8")
		return t.Execute(ctx.ResponseWriter, data)
	} else {
		return fmt.Errorf("template %q is undefined", name)
	}
}
