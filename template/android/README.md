# TGE-CLI Android Target 
## Requirements
 * GO 1.12
 * Android NDK R19+ - https://developer.android.com/ndk/guides

## Command line
```shell
tge-cli build -target android [package-path]
```
No specific options available for Android but the -dev flag allows to build an APK for all Android targets architectures.

## Customization
The app icon can be modified using the icon.png file.

The Manifest can also be edited to set app name or screen orientation.

## Release content
With -dev set (debug) :
```
/dist
 |-- android
     |-- ${APP-NAME}.apk
```

Without set (release) :
```
/dist
|-- android
    |-- ${APP-NAME}-386.apk
    |-- ${APP-NAME}-amd64.apk
    |-- ${APP-NAME}-arm.apk
    |-- ${APP-NAME}-arm64.apk
```

Assets are packaged in APK. 

Signature options are not implemented yet but planned.