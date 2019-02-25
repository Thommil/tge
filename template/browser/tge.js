(() => {
    if (typeof window !== 'undefined') {
        window.global = window;
    } else if (typeof self !== 'undefined') {
        self.global = self;
    } else {
        throw new Error('cannot start TGE (neither window nor self is defined)');
    }

    let canvasEl = document.getElementById('canvas');

    if (!canvasEl) {
        throw new Error('Canvas element not found (must be #canvas)');
    }

    let assetsMap = {}

    global.tge = {
        init() {
            canvasEl.classList.remove('stop');
            canvasEl.classList.add('start');
            canvasEl.oncontextmenu = function (e) {e.preventDefault();};
            canvasEl.focus()
            return canvasEl;
        },

        setFullscreen(enabled) {
            if (enabled) {
                canvasEl.classList.add('fullscreen');
            } else {
                canvasEl.classList.remove('fullscreen');
            }
            canvasEl.setAttribute('width',canvasEl.clientWidth);
            canvasEl.setAttribute('height', canvasEl.clientHeight);
        },

        resize(width, height) {
            canvasEl.style['width'] = width + 'px';
            canvasEl.style['height'] = height + 'px';
            canvasEl.setAttribute('width',canvasEl.clientWidth);
            canvasEl.setAttribute('height', canvasEl.clientHeight);
        },

        getAssetSize(path, callback) {            
            fetch('./assets/' + path).then((response) => {
                if(response.ok) {                    
                    return response.arrayBuffer()           
                } else {
                    throw new Error(response.statusText)
                }
            })
            .then((content) => {
                if (content) {
                    assetsMap[path] = new Uint8Array(content)                    
                } else {
                    throw new Error("empty content") 
                }
                callback(content.byteLength, null)
            })
            .catch((error) => {
                callback(null, error)
            });
        },

        loadAsset(path, goData, callback) {
            if (assetsMap[path]) {
                goData.set(assetsMap[path])
                delete assetsMap[path]
                callback(null)
            } else {
                callback("empty content")
            }
        },

        stop() {
            canvasEl.classList.remove('start');
            canvasEl.classList.add('stop');
        },

        showError(err) {
            console.error(err);
        }
    }

    window.onload = function(){
        window.go = new Go();
        if(WebAssembly.instantiateStreaming) {
            WebAssembly.instantiateStreaming(fetch("main.wasm"), window.go.importObject).then((result) => {
                window.go.run(result.instance);
            }).catch(err=>{
                tge.showError(err)
            });
        } else {
            fetch('main.wasm').then(response =>
                response.arrayBuffer()
              ).then(bytes =>
                WebAssembly.instantiate(bytes, window.go.importObject)
              ).then(result => {
                window.go.run(result.instance);
            }).catch(err=>{
                tge.showError(err)
            });
        }
    }
})();