htmx.defineExtension('format-response', {
    onEvent: function(name, evt) {
        if (name === "htmx:beforeOnLoad") {
            let raw = evt.detail.xhr.responseText;
            evt.detail.serverResponse = formatContent(raw);
        }
}});

function formatContent(html) {
    return `<div class="formatted">${html}</div>`;
}

