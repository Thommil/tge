# TGE - Game Runtime in GO
**TGE** aims to provide a light, portable and almsot unopiniated runtime to integrate your favorite GO libraries (at least mines).

## TGE Core 
Core implements a minimal runtime for several platforms:
  * Windows, MacOS & Linux
  * Android 5+ & IOS 8+
  * Web on Chrome, Firefox & Safari

For diagrams lovers, here's the main architecture:

![Core Architecture](https://raw.githubusercontent.com/thommil/tge/master/specs/api.png)

## TGE plugins
Plugins then allow to create your App by choosing implementation of different parts (GUI, rendering, sound, physics ...). If a library needs a custom implementation to run on each target, I try to create the associated plugin, if not, just import it by yourself. This way, **you choose** how to create your game/application. 

Currently available plugins :
 * [tge-gl](https://github.com/thommil/tge-gl) : OpenGL/ES 3+ API 
 * [tge-g3n](https://github.com/thommil/tge-g3n) : Port of the awesome [G3N](https://github.com/g3n/engine) Game Engine




# Getting started
Based on

