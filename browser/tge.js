// TGE Tooling JS
(() => {
    if (typeof window !== "undefined") {
        window.global = window;
    } else if (typeof self !== "undefined") {
        self.global = self;
    } else {
        throw new Error("cannot start TGE (neither window nor self is defined)");
    }

    
})();