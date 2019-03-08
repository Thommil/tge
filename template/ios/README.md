# TGE-CLI IOS Target 
## Requirements
 * GO 1.12
 * MacOS 10+
 * Apple Developer Account - http://developer.apple.com
 * Existing Provisioning Profile for the App

## Command line
```shell
tge-cli build -target ios -bundleid BUNDLE_ID [package-path]
```
The -bundleid flag is mandatory for IOS build, this ID can be obtained from Apple Developer Console or XCode.

## Customization
The app icon can be modified using the icon.png file.

## Release content
... WIP