# TGE-CLI Browser Target 
## Requirements
 * GO 1.12

## Command line
```shell
tge-cli build -target browser [package-path]
```

## Customization
The favicon can be modified using the favicon.ico file.

The tge-min.js is manualy minified, future releases will use Webpack for this task and PWA support.

Feel free to customize the HTML/CSS/JS code to cover your needs, the canvas as been implemented that way.

## Release content
The following content must be uploaded to your Web server:

```
/dist
|-- /browser
    |-- /assets
    |-- favicon.ico
    |-- index.html
    |-- tge-min.js
    |-- tge.css
    |-- main.wasm
```

Use [goexec](https://github.com/shurcooL/goexec) to test:
```shell
goexec 'http.ListenAndServe(":8080", http.FileServer(http.Dir("/path/to/release")))'
```

It's strongly advised to enable gzip on your web server as the size of the WASM file can be drastically by compression. 