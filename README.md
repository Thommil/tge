# TGE - Game Runtime in GO
**TGE** aims to provide a light, portable and almsot unopiniated runtime to integrate your favorite GO libraries (at least mine).

**TGE Core** implements a minimal runtime for several platforms:
  * Windows, MacOS & Linux
  * Android 5+ & IOS 8+
  * Web on Chrome, Firefox & Safari



```plantuml
@startuml API
' Components
interface App {
    OnCreate(Settings) error
    OnStart(Runtime) error
    OnResize(int, int)
    OnResume()
    OnRender(Duration, Mutex)
    OnTick(Duration, Mutex)
    OnMouseEvent(MouseEvent)
    OnScrollEvent(ScrollEvent)
    OnKeyEvent(KeyEvent)
    OnPause()
    OnStop()
    OnDispose()
}

class tge << (P,#FF7700) Package >> {
    {static} Run(App)
}

interface Runtime {
    Use(Plugin)
    GetPlugin(string) []Plugin
    GetHost() interface{}
    GetRenderer() interface{}
    LoadAsset(string) []byte, error
    Stop()
}

interface Plugin{    
    Init(Runtime) error
    GetName() string
    Dispose()
}

class "tge-*" << (P,#FF7700) Packages >> {
    GetPlugin() Plugin
}

' Relations
App --> tge : Run(App)
tge --> Runtime : instanciate
App <-- Runtime : manage
App --> Runtime : use
App --> Plugin : use
Runtime --> Plugin : manage
Plugin --> Runtime  : use
"tge-*" --> Plugin : expose
App ..> "tge-*" : load

@enduml
```

# Getting started
Based on

