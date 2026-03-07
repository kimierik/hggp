htmx.defineExtension('format-response', {
    onEvent: function(name, evt) {
        if (name === "htmx:beforeOnLoad") {
            let raw = evt.detail.xhr.responseText;
            evt.detail.serverResponse = formatContent(raw);
        }
}});

function formatContent(obj) {
    obj = JSON.parse(obj)
    return `
        <div class="formatted-message border p-1">
            <p class="border-bottom">name: ${obj.author}</p>
            <p>${obj.message}</p>
        </div>`;
}
