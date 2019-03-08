# TGE-CLI Windows Target 
## Requirements
 * GO 1.12
 * A cgo compiler (gcc)

A simple solution is to use Git Bash (https://git-scm.com/downloads) with MinGW support (http://mingw.org/).

## Command line
```shell
tge-cli build -target desktop [package-path]
```

## Customization
The icon can be modified using the icon.png file.

The verisioninfo.json file allow to edit applications attibutes displayed in program properties of Windows (name, version, description...).

## Release content
With -dev set (debug) :
```
/dist
|-- /windows
    |-- ${APP-NAME}.exe
```
In dev mode, assets are retrieved from package path and not copied at each build. The dev mode just create an executable and allow to see console ouput.

Without -dev set (release) :
```
/dist
|-- /windows
    |-- /assets
    |-- ${APP-NAME}.exe
```
A true Windows application is created in this mode, assets are in a dedicated folder copied from package assets.