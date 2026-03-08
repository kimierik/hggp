
# concept
This is more of a test/learning thing for these technologies so this is not an actual project now.
More of a proof of concept and a thing that lets me know about the different usecases.
Completely unfinished and unpolished and i will probably not make attempts to fix it, since it has served its purpoise to me.

# dependencies

i know it is a bit mutch....
```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
google.golang.org/genproto/googleapis/api
google.golang.org/genproto/googleapis/rpc
google.golang.org/grpc v1.79.1
google.golang.org/protobuf v1.36.11
```

# running
this is a bit of a mess (partly on purpoise since the idea was never to make this into a real thing) but you must run ./build.sh and the commands in ./backend/protoc

then 
```sh
docker compose run --build --detach
```


# notes/ideas/misc ramble

with wasm you can use and instanciate go templates from the client. 
So you can use grpc endpoints or any other endpoint that returns non html data. 
This breaks HATEOS princibles that HTMX loves but some apps might need to query data from other places or need to be a bit more client heavy.
In this case if you have a multipurpoise backend where you use grpc then you might either need to make an endpoint for for each thing that the client wants to touch just so that you can wrap it in html, this would be alot of work.
  

With the examlple that i have here. you can use HTMX to send requests and use extensions to format them without having to have the rendering logic be exclusively in js.
in index.html the 'format-response' extention just formats a json object into html. Ideally you could extend this to format the json into whatever go template with a wasm call.
rn im too lazy to do it so im just going to explain it.
```js
async function f(){
    await window.GoReady
    let a =App.renderTemplate("foba", {"item":"another", "foobar":124})

    if (a instanceof Error){
        throw(a)
    }
    return a
};f()
```
You would create an api like this that would render the template 'foba', with the following arguments, 
in an hx request you can have the name of the template in html attributes, so the rendering logic is very explicit in the html.

you would have to remake the frontend/static/js/htmx_extention.js into somthng that looks more like this.
```js
htmx.defineExtension('format-response', {
    onEvent: function(name, evt) {
        if (name === "htmx:beforeOnLoad") {
            let raw = evt.detail.xhr.responseText;
            let templateName = getTemplateNameFromEvt(evt)
            evt.detail.serverResponse = App.renderTemplate(templateName, JSON.parse(raw))
        }
}});
```

  
  
Then in wasm.go you could have somthing like this. 
```go
func renderTemplate(this js.Value, args []js.Value) interface{} {
	var tmpl = template.Must(
		template.ParseFS(templatesFS, args[0]),
	)
    // rest of the code for rendering the template
```
This gives you the power of having the templates at the server so you can return server rendered html and use the same templates at the client.
The code for rendering the templates is also very simple and can mostly be expressed in html and not seperately in a js file somewhere.


# TLDR
- Server and client can use the same go templates for rendering.
- Logic for how the app behaves can be moved to HTML with HTMX and HTMX-extensions. This way you do not have important rendering/hydrating logic in a seperate js file.
- Singular static binary for backend and frontend with Go.
- Building can be made simple with a makefile.
- gRPC is cool and can be used to make usefull endpoints that can be interacted with

# things left unimplemented
- Auth
- build system

# Things that can still be optimized/improved
- Generating SQL definitions from protobuf, or having any orm/generated database layer.
