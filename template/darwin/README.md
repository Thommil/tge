# TGE-CLI MacOS Target 
## Requirements
 * GO 1.12
 * A cgo compiler (gcc)

## Command line
```shell
tge-cli build -target desktop [package-path]
```

## Customization
The icon can be modified using the icon.png file.

## Release content
With -dev set (debug) :
```
/dist
|-- /darwin
    |-- ${APP-NAME}
```
In dev mode, assets are retrieved from package path and not copied at each build. The dev mode just create an executable and allow to see console ouput.

Without -dev set (release) :
```
/dist
|-- /darwin
    |-- ${APP-NAME}.app
```
A true MacOS application is created in this mode, assets are embedded in the .app package.