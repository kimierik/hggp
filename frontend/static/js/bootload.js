window.GoReady = (async () => {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(
        fetch("/static/wasm/app.wasm"),
        go.importObject
    );
    go.run(result.instance);
})();

/*
async function f(){
    await window.GoReady
    let a =App.renderTemplate("foba", {"item":"another", "foobar":124})

    //let a =App.renderTemplate("fobba")
    if (a instanceof Error){
        console.log("render template",a)
        throw(a)
    }
}
f()
 * */
