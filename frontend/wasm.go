package main

import (
	"embed"
	"html/template"
	"bytes"
	"fmt"
	"syscall/js"
)

//go:embed templates/*.html
var templatesFS embed.FS




func renderTemplate(this js.Value, args []js.Value) interface{} {
	var tmpl = template.Must(
		template.ParseFS(templatesFS, "templates/header.html"),
	)
	if len(args)<=1{
		err := js.Global().Get("Error").New("No arguments provided")
        // Panic with the JS Error
		return err

	}

	for i,arg:=range args{
		fmt.Println(i,arg)
	}

	var buf bytes.Buffer
	tmpl.ExecuteTemplate(&buf, "header", nil)

	return buf.String()
}

func main() {

	app := js.Global().Get("Object").New()

	// function decls here
	app.Set("renderTemplate", js.FuncOf(renderTemplate))

	js.Global().Set("App", app)
	select {}
}
