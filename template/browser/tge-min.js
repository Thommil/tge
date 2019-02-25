//Fetch Polyfill from https://github.github.io/fetch/
!function(t,e){"object"==typeof exports&&"undefined"!=typeof module?e(exports):"function"==typeof define&&define.amd?define(["exports"],e):e(t.WHATWGFetch={})}(this,function(a){"use strict";var r="URLSearchParams"in self,o="Symbol"in self&&"iterator"in Symbol,h="FileReader"in self&&"Blob"in self&&function(){try{return new Blob,!0}catch(t){return!1}}(),n="FormData"in self,i="ArrayBuffer"in self;if(i)var e=["[object Int8Array]","[object Uint8Array]","[object Uint8ClampedArray]","[object Int16Array]","[object Uint16Array]","[object Int32Array]","[object Uint32Array]","[object Float32Array]","[object Float64Array]"],s=ArrayBuffer.isView||function(t){return t&&-1<e.indexOf(Object.prototype.toString.call(t))};function u(t){if("string"!=typeof t&&(t=String(t)),/[^a-z0-9\-#$%&'*+.^_`|~]/i.test(t))throw new TypeError("Invalid character in header field name");return t.toLowerCase()}function f(t){return"string"!=typeof t&&(t=String(t)),t}function t(e){var t={next:function(){var t=e.shift();return{done:void 0===t,value:t}}};return o&&(t[Symbol.iterator]=function(){return t}),t}function d(e){this.map={},e instanceof d?e.forEach(function(t,e){this.append(e,t)},this):Array.isArray(e)?e.forEach(function(t){this.append(t[0],t[1])},this):e&&Object.getOwnPropertyNames(e).forEach(function(t){this.append(t,e[t])},this)}function c(t){if(t.bodyUsed)return Promise.reject(new TypeError("Already read"));t.bodyUsed=!0}function p(r){return new Promise(function(t,e){r.onload=function(){t(r.result)},r.onerror=function(){e(r.error)}})}function y(t){var e=new FileReader,r=p(e);return e.readAsArrayBuffer(t),r}function l(t){if(t.slice)return t.slice(0);var e=new Uint8Array(t.byteLength);return e.set(new Uint8Array(t)),e.buffer}function b(){return this.bodyUsed=!1,this._initBody=function(t){var e;(this._bodyInit=t)?"string"==typeof t?this._bodyText=t:h&&Blob.prototype.isPrototypeOf(t)?this._bodyBlob=t:n&&FormData.prototype.isPrototypeOf(t)?this._bodyFormData=t:r&&URLSearchParams.prototype.isPrototypeOf(t)?this._bodyText=t.toString():i&&h&&((e=t)&&DataView.prototype.isPrototypeOf(e))?(this._bodyArrayBuffer=l(t.buffer),this._bodyInit=new Blob([this._bodyArrayBuffer])):i&&(ArrayBuffer.prototype.isPrototypeOf(t)||s(t))?this._bodyArrayBuffer=l(t):this._bodyText=t=Object.prototype.toString.call(t):this._bodyText="",this.headers.get("content-type")||("string"==typeof t?this.headers.set("content-type","text/plain;charset=UTF-8"):this._bodyBlob&&this._bodyBlob.type?this.headers.set("content-type",this._bodyBlob.type):r&&URLSearchParams.prototype.isPrototypeOf(t)&&this.headers.set("content-type","application/x-www-form-urlencoded;charset=UTF-8"))},h&&(this.blob=function(){var t=c(this);if(t)return t;if(this._bodyBlob)return Promise.resolve(this._bodyBlob);if(this._bodyArrayBuffer)return Promise.resolve(new Blob([this._bodyArrayBuffer]));if(this._bodyFormData)throw new Error("could not read FormData body as blob");return Promise.resolve(new Blob([this._bodyText]))},this.arrayBuffer=function(){return this._bodyArrayBuffer?c(this)||Promise.resolve(this._bodyArrayBuffer):this.blob().then(y)}),this.text=function(){var t,e,r,o=c(this);if(o)return o;if(this._bodyBlob)return t=this._bodyBlob,e=new FileReader,r=p(e),e.readAsText(t),r;if(this._bodyArrayBuffer)return Promise.resolve(function(t){for(var e=new Uint8Array(t),r=new Array(e.length),o=0;o<e.length;o++)r[o]=String.fromCharCode(e[o]);return r.join("")}(this._bodyArrayBuffer));if(this._bodyFormData)throw new Error("could not read FormData body as text");return Promise.resolve(this._bodyText)},n&&(this.formData=function(){return this.text().then(v)}),this.json=function(){return this.text().then(JSON.parse)},this}d.prototype.append=function(t,e){t=u(t),e=f(e);var r=this.map[t];this.map[t]=r?r+", "+e:e},d.prototype.delete=function(t){delete this.map[u(t)]},d.prototype.get=function(t){return t=u(t),this.has(t)?this.map[t]:null},d.prototype.has=function(t){return this.map.hasOwnProperty(u(t))},d.prototype.set=function(t,e){this.map[u(t)]=f(e)},d.prototype.forEach=function(t,e){for(var r in this.map)this.map.hasOwnProperty(r)&&t.call(e,this.map[r],r,this)},d.prototype.keys=function(){var r=[];return this.forEach(function(t,e){r.push(e)}),t(r)},d.prototype.values=function(){var e=[];return this.forEach(function(t){e.push(t)}),t(e)},d.prototype.entries=function(){var r=[];return this.forEach(function(t,e){r.push([e,t])}),t(r)},o&&(d.prototype[Symbol.iterator]=d.prototype.entries);var m=["DELETE","GET","HEAD","OPTIONS","POST","PUT"];function w(t,e){var r,o,n=(e=e||{}).body;if(t instanceof w){if(t.bodyUsed)throw new TypeError("Already read");this.url=t.url,this.credentials=t.credentials,e.headers||(this.headers=new d(t.headers)),this.method=t.method,this.mode=t.mode,this.signal=t.signal,n||null==t._bodyInit||(n=t._bodyInit,t.bodyUsed=!0)}else this.url=String(t);if(this.credentials=e.credentials||this.credentials||"same-origin",!e.headers&&this.headers||(this.headers=new d(e.headers)),this.method=(r=e.method||this.method||"GET",o=r.toUpperCase(),-1<m.indexOf(o)?o:r),this.mode=e.mode||this.mode||null,this.signal=e.signal||this.signal,this.referrer=null,("GET"===this.method||"HEAD"===this.method)&&n)throw new TypeError("Body not allowed for GET or HEAD requests");this._initBody(n)}function v(t){var n=new FormData;return t.trim().split("&").forEach(function(t){if(t){var e=t.split("="),r=e.shift().replace(/\+/g," "),o=e.join("=").replace(/\+/g," ");n.append(decodeURIComponent(r),decodeURIComponent(o))}}),n}function E(t,e){e||(e={}),this.type="default",this.status=void 0===e.status?200:e.status,this.ok=200<=this.status&&this.status<300,this.statusText="statusText"in e?e.statusText:"OK",this.headers=new d(e.headers),this.url=e.url||"",this._initBody(t)}w.prototype.clone=function(){return new w(this,{body:this._bodyInit})},b.call(w.prototype),b.call(E.prototype),E.prototype.clone=function(){return new E(this._bodyInit,{status:this.status,statusText:this.statusText,headers:new d(this.headers),url:this.url})},E.error=function(){var t=new E(null,{status:0,statusText:""});return t.type="error",t};var A=[301,302,303,307,308];E.redirect=function(t,e){if(-1===A.indexOf(e))throw new RangeError("Invalid status code");return new E(null,{status:e,headers:{location:t}})},a.DOMException=self.DOMException;try{new a.DOMException}catch(t){a.DOMException=function(t,e){this.message=t,this.name=e;var r=Error(t);this.stack=r.stack},a.DOMException.prototype=Object.create(Error.prototype),a.DOMException.prototype.constructor=a.DOMException}function _(n,s){return new Promise(function(o,t){var e=new w(n,s);if(e.signal&&e.signal.aborted)return t(new a.DOMException("Aborted","AbortError"));var i=new XMLHttpRequest;function r(){i.abort()}i.onload=function(){var t,n,e={status:i.status,statusText:i.statusText,headers:(t=i.getAllResponseHeaders()||"",n=new d,t.replace(/\r?\n[\t ]+/g," ").split(/\r?\n/).forEach(function(t){var e=t.split(":"),r=e.shift().trim();if(r){var o=e.join(":").trim();n.append(r,o)}}),n)};e.url="responseURL"in i?i.responseURL:e.headers.get("X-Request-URL");var r="response"in i?i.response:i.responseText;o(new E(r,e))},i.onerror=function(){t(new TypeError("Network request failed"))},i.ontimeout=function(){t(new TypeError("Network request failed"))},i.onabort=function(){t(new a.DOMException("Aborted","AbortError"))},i.open(e.method,e.url,!0),"include"===e.credentials?i.withCredentials=!0:"omit"===e.credentials&&(i.withCredentials=!1),"responseType"in i&&h&&(i.responseType="blob"),e.headers.forEach(function(t,e){i.setRequestHeader(e,t)}),e.signal&&(e.signal.addEventListener("abort",r),i.onreadystatechange=function(){4===i.readyState&&e.signal.removeEventListener("abort",r)}),i.send(void 0===e._bodyInit?null:e._bodyInit)})}_.polyfill=!0,self.fetch||(self.fetch=_,self.Headers=d,self.Request=w,self.Response=E),a.Headers=d,a.Request=w,a.Response=E,a.fetch=_,Object.defineProperty(a,"__esModule",{value:!0})});
//WASM
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
(()=>{if("undefined"!=typeof window)window.global=window;else if("undefined"!=typeof self)self.global=self;else throw new Error("cannot export Go (neither window nor self is defined)");let a="";global.fs={constants:{O_WRONLY:-1,O_RDWR:-1,O_CREAT:-1,O_TRUNC:-1,O_APPEND:-1,O_EXCL:-1},writeSync(b,d){a+=c.decode(d);const e=a.lastIndexOf("\n");return-1!=e&&(console.log(a.substr(0,e)),a=a.substr(e+1)),d.length},openSync(){const a=new Error("not implemented");throw a.code="ENOSYS",a}};const b=new TextEncoder("utf-8"),c=new TextDecoder("utf-8");global.Go=class{constructor(){this.argv=["js"],this.env={},this.exit=a=>{0!==a&&console.warn("exit code:",a)},this._callbackTimeouts=new Map,this._nextCallbackTimeoutID=1;const a=()=>new DataView(this._inst.exports.mem.buffer),d=(b,c)=>{a().setUint32(b+0,c,!0),a().setUint32(b+4,Math.floor(c/4294967296),!0)},e=b=>{const c=a().getUint32(b+0,!0),d=a().getInt32(b+4,!0);return c+4294967296*d},f=b=>{const c=a().getFloat64(b,!0);if(!isNaN(c))return c;const d=a().getUint32(b,!0);return this._values[d]},g=(b,c)=>{const d=2146959360;if("number"==typeof c)return isNaN(c)?(a().setUint32(b+4,2146959360,!0),void a().setUint32(b,0,!0)):void a().setFloat64(b,c,!0);switch(c){case void 0:return a().setUint32(b+4,d,!0),void a().setUint32(b,1,!0);case null:return a().setUint32(b+4,d,!0),void a().setUint32(b,2,!0);case!0:return a().setUint32(b+4,d,!0),void a().setUint32(b,3,!0);case!1:return a().setUint32(b+4,d,!0),void a().setUint32(b,4,!0);}let e=this._refs.get(c);e===void 0&&(e=this._values.length,this._values.push(c),this._refs.set(c,e));let f=0;switch(typeof c){case"string":f=1;break;case"symbol":f=2;break;case"function":f=3;}a().setUint32(b+4,2146959360|f,!0),a().setUint32(b,e,!0)},h=a=>{const b=e(a+0),c=e(a+8);return new Uint8Array(this._inst.exports.mem.buffer,b,c)},i=b=>{const c=e(b+0),d=e(b+8),g=Array(d);for(let a=0;a<d;a++)g[a]=f(c+8*a);return g},j=a=>{const b=e(a+0),d=e(a+8);return c.decode(new DataView(this._inst.exports.mem.buffer,b,d))},k=Date.now()-performance.now();this.importObject={go:{"runtime.wasmExit":b=>{const c=a().getInt32(b+8,!0);this.exited=!0,delete this._inst,delete this._values,delete this._refs,this.exit(c)},"runtime.wasmWrite":b=>{const c=e(b+8),d=e(b+16),f=a().getInt32(b+24,!0);fs.writeSync(c,new Uint8Array(this._inst.exports.mem.buffer,d,f))},"runtime.nanotime":a=>{d(a+8,1e6*(k+performance.now()))},"runtime.walltime":b=>{const c=new Date().getTime();d(b+8,c/1e3),a().setInt32(b+16,1e6*(c%1e3),!0)},"runtime.scheduleCallback":b=>{const c=this._nextCallbackTimeoutID;this._nextCallbackTimeoutID++,this._callbackTimeouts.set(c,setTimeout(()=>{this._resolveCallbackPromise()},e(b+8)+1)),a().setInt32(b+16,c,!0)},"runtime.clearScheduledCallback":b=>{const c=a().getInt32(b+8,!0);clearTimeout(this._callbackTimeouts.get(c)),this._callbackTimeouts.delete(c)},"runtime.getRandomData":a=>{crypto.getRandomValues(h(a+8))},"syscall/js.stringVal":a=>{g(a+24,j(a+8))},"syscall/js.valueGet":a=>{g(a+32,Reflect.get(f(a+8),j(a+16)))},"syscall/js.valueSet":a=>{Reflect.set(f(a+8),j(a+16),f(a+32))},"syscall/js.valueIndex":a=>{g(a+24,Reflect.get(f(a+8),e(a+16)))},"syscall/js.valueSetIndex":a=>{Reflect.set(f(a+8),e(a+16),f(a+24))},"syscall/js.valueCall":b=>{try{const c=f(b+8),d=Reflect.get(c,j(b+16)),e=i(b+32);g(b+56,Reflect.apply(d,c,e)),a().setUint8(b+64,1)}catch(c){g(b+56,c),a().setUint8(b+64,0)}},"syscall/js.valueInvoke":b=>{try{const c=f(b+8),d=i(b+16);g(b+40,Reflect.apply(c,void 0,d)),a().setUint8(b+48,1)}catch(c){g(b+40,c),a().setUint8(b+48,0)}},"syscall/js.valueNew":b=>{try{const c=f(b+8),d=i(b+16);g(b+40,Reflect.construct(c,d)),a().setUint8(b+48,1)}catch(c){g(b+40,c),a().setUint8(b+48,0)}},"syscall/js.valueLength":a=>{d(a+16,parseInt(f(a+8).length))},"syscall/js.valuePrepareString":a=>{const c=b.encode(f(a+8)+"");g(a+16,c),d(a+24,c.length)},"syscall/js.valueLoadString":a=>{const b=f(a+8);h(a+16).set(b)},"syscall/js.valueInstanceOf":b=>{a().setUint8(b+24,f(b+8)instanceof f(b+16))},debug:a=>{console.log(a)}}}}async run(a){this._inst=a,this._values=[NaN,void 0,null,!0,!1,global,this._inst.exports.mem,this],this._refs=new Map,this._callbackShutdown=!1,this.exited=!1;const c=new DataView(this._inst.exports.mem.buffer);let d=4096;const e=a=>{let e=d;return new Uint8Array(c.buffer,d,a.length+1).set(b.encode(a+"\0")),d+=a.length+(8-a.length%8),e},f=this.argv.length,g=[];this.argv.forEach(a=>{g.push(e(a))});const h=Object.keys(this.env).sort();g.push(h.length),h.forEach(a=>{g.push(e(`${a}=${this.env[a]}`))});const i=d;for(g.forEach(a=>{c.setUint32(d,a,!0),c.setUint32(d+4,0,!0),d+=8});;){const a=new Promise(a=>{this._resolveCallbackPromise=()=>{if(this.exited)throw new Error("bad callback: Go program has already exited");setTimeout(a,0)}});if(this._inst.exports.run(f,i),this.exited)break;await a}}static _makeCallbackHelper(a,b,c){return function(){b.push({id:a,args:arguments}),c._resolveCallbackPromise()}}static _makeEventCallbackHelper(a,b,c,d){return function(e){a&&e.preventDefault(),b&&e.stopPropagation(),c&&e.stopImmediatePropagation(),d(e)}}}})();
// TGE Tooling JS
(()=>{if("undefined"!=typeof window)window.global=window;else if("undefined"!=typeof self)self.global=self;else throw new Error("cannot start TGE (neither window nor self is defined)");let a=document.getElementById("canvas");if(!a)throw new Error("Canvas element not found (must be #canvas)");let b={};global.tge={init(){return a.classList.remove("stop"),a.classList.add("start"),a.oncontextmenu=function(a){a.preventDefault()},a.focus(),a},setFullscreen(b){b?a.classList.add("fullscreen"):a.classList.remove("fullscreen"),a.setAttribute("width",a.clientWidth),a.setAttribute("height",a.clientHeight)},resize(b,c){a.style.width=b+"px",a.style.height=c+"px",a.setAttribute("width",a.clientWidth),a.setAttribute("height",a.clientHeight)},getAssetSize(a,c){fetch("./assets/"+a).then(a=>{if(a.ok)return a.arrayBuffer();throw new Error(a.statusText)}).then(d=>{if(d)b[a]=new Uint8Array(d);else throw new Error("empty content");c(d.byteLength,null)}).catch(a=>{c(null,a)})},loadAsset(a,c,d){b[a]?(c.set(b[a]),delete b[a],d(null)):d("empty content")},stop(){a.classList.remove("start"),a.classList.add("stop")},showError(a){console.error(a)}},window.onload=function(){window.go=new Go,WebAssembly.instantiateStreaming?WebAssembly.instantiateStreaming(fetch("main.wasm"),window.go.importObject).then(a=>{window.go.run(a.instance)}).catch(a=>{tge.showError(a)}):fetch("main.wasm").then(a=>a.arrayBuffer()).then(a=>WebAssembly.instantiate(a,window.go.importObject)).then(a=>{window.go.run(a.instance)}).catch(a=>{tge.showError(a)})}})();