#!/bin/sh

# Env/Opts
usage() { echo "Usage: build.sh [-t <desktop|browser|android|ios>] PROJECT_PATH" 1>&2; exit 1; }

while getopts "t:" o; do
    case "${o}" in
        t)
            TARGET=${OPTARG}
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

if [ -z "$TARGET" ]; then
    usage
fi

PROJECT_PATH="$(cd "$(dirname "$@")"; pwd)/$(basename "$@")"

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
    echo "ERROR : No implemented yet"
    exit 3
}

# Build IOS
BuildIOS () {
    echo "ERROR : No implemented yet"
    exit 3
}

# Build browser
BuildBrowser () {
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
}
    

# Build
if [ "$TARGET" == "desktop" ] ; then
    BuildDesktop
elif [ "$TARGET" == "android" ] ; then
    BuildAndroid
elif [ "$TARGET" == "ios" ] ; then
    BuildIOS
elif [ "$TARGET" == "browser" ] ; then
    BuildBrowser
else
    echo "ERROR : Unsupported target : $TARGET"
fi

