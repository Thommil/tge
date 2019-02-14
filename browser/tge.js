// TGE Tooling JS
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

    global.tge = {
        setFullscreen(enabled) {
            if (enabled) {
                canvasEl.classList.add('fullscreen');
            } else {
                canvasEl.classList.remove('fullscreen');
            }
        },

        resize(width, height) {
            canvasEl.style['width'] = width + 'px';
            canvasEl.style['height'] = height + 'px';
        },

        init() {
            canvasEl.classList.remove('stop');
            canvasEl.classList.add('start');
            return canvasEl;
        },

        stop() {
            canvasEl.classList.remove('start');
            canvasEl.classList.add('stop');
        }

    }
})();