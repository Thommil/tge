#!/bin/sh

# Env
if [ $# -lt 2 ] ; then
    echo "Synopsis : build.sh PROJECT_PATH TARGET [TGE_BRANCH]"
    echo "  Available targets  : desktop android ios browser"
    exit 1
fi

if [ -z "$GOPATH" ] ; then
    echo "ERROR : GOPATH is not set"
    exit 1
fi 

PROJECT_PATH="$(cd "$(dirname "$1")"; pwd)/$(basename "$1")"
TARGET="$2"

if [ ! -d "$PROJECT_PATH" ] ; then
    echo "ERROR : Project path not found : $PROJECT_PATH"
    exit 2
fi

if [ "$(uname)" == "Darwin" ]; then
    OS="macos"
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    OS="linux"
elif [ "$(expr substr $(uname -s) 1 5)" == "MINGW" ]; then
    OS="windows"
else
    echo "ERROR : Unsupported OS : $(uname -s)"
    exit 3
fi

# Init
BUILDER_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DIST_PATH="$PROJECT_PATH/dist/$TARGET"

# Build desktop
BuildDesktop () {
    echo "Building $PROJECT_PATH for desktop ..."  
    rm -rf "$DIST_PATH"
    mkdir -p "$DIST_PATH"
    cd $PROJECT_PATH >/dev/null
    go build -o "$DIST_PATH/$(basename $PROJECT_PATH)"    
    if [ "$?" -eq "0" ]; then
        if [ $OS == "windows" ]; then
            mv "$DIST_PATH/$(basename $PROJECT_PATH)" "$DIST_PATH/$(basename $PROJECT_PATH).exe"
        fi
        echo "Build success in $DIST_PATH"  
    else
        rm -rf "$DIST_PATH"
    fi
    cd - > /dev/null
}

# Build Android
BuildAndroid () {
    # echo "Building $PROJECT_PATH for android ..."
    echo "ERROR : No implemented yet"
    exit 3
}

# Build IOS
BuildIOS () {
    # echo "Building $PROJECT_PATH for ios ..."
    echo "ERROR : No implemented yet"
    exit 3
}

# Build browser
BuildBrowser () {
    echo "Building $PROJECT_PATH for browser ..."
    rm -rf "$DIST_PATH"
    mkdir -p "$DIST_PATH"
    cd $PROJECT_PATH >/dev/null
    GOOS=js GOARCH=wasm go build -o "$DIST_PATH/main.wasm"    
    if [ "$?" -eq "0" ]; then
        cp $BUILDER_PATH/browser/index.html $DIST_PATH
        cp $BUILDER_PATH/browser/tge.js $DIST_PATH
        cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" $DIST_PATH
        echo "Build success in $DIST_PATH"  
    else
        rm -rf "$DIST_PATH"
    fi
    cd - > /dev/null
    # cd $PROJECTS_PATH/$PROJECT
    # GOARCH=js GOOS=browser go build
    # cd -
    # cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
    # echo "ERROR : No implemented yet"
    # exit 3
}
    

# Build
if [ $TARGET == "desktop" ] ; then
    BuildDesktop
elif [ $TARGET == "android" ] ; then
    BuildAndroid
elif [ $TARGET == "ios" ] ; then
    BuildIOS
elif [ $TARGET == "browser" ] ; then
    BuildBrowser
else
    echo "ERROR : Unsupported target : $TARGET"
fi

