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
    echo "ERROR : Project path not found : $PROJECT_PATH" 1>&2
    exit 2
fi

if [ "$(uname)" == "Darwin" ]; then
    OS="macos"
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    OS="linux"
elif [ "$(expr substr $(uname -s) 1 5)" == "MINGW" ]; then
    OS="windows"
else
    echo "ERROR : Unsupported OS : $(uname -s)" 1>&2
    exit 3
fi

# Init
BUILDER_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DIST_PATH="$PROJECT_PATH/dist/$TARGET"

# Build desktop
BuildDesktop () {
    echo " > cleaning"
    rm -rf "$DIST_PATH"
    mkdir -p "$DIST_PATH"
    cd $PROJECT_PATH >/dev/null

    echo " > building"
    go build -o "$DIST_PATH/$(basename $PROJECT_PATH)"    

    if [ "$?" -eq "0" ]; then
        if [ $OS == "windows" ]; then
            mv "$DIST_PATH/$(basename $PROJECT_PATH)" "$DIST_PATH/$(basename $PROJECT_PATH).exe"
        fi

        echo " > copying resources"
        if [ ! -d "$PROJECT_PATH/desktop" ]; then
            cp -rp "$BUILDER_PATH/desktop" "$PROJECT_PATH"
        fi
        cp -rp "$PROJECT_PATH/desktop" $(dirname $DIST_PATH)
        
        echo " > Build done in $DIST_PATH"  
    else
        echo " > Build failed" 1>&2
        rm -rf "$DIST_PATH"
    fi

    cd - > /dev/null
}

# Build Android
BuildAndroid () {
    echo " > init gomobile"
    if [ "$ANDROID_NDK" == "" ]; then
        echo "ERROR : ANDROID_NDK is not set (should be \$ANDROID_HOME/ndk-bundle)" >&2
        echo " > Build failed" 1>&2
        exit 5
    fi
    command -v gomobile >/dev/null 2>&1 || { echo "ERROR : gomobile command not found in PATH" 1>&2; exit 1; }
    #gomobile init -ndk $ANDROID_NDK
    
    echo " > cleaning"
    rm -rf "$DIST_PATH"
    mkdir -p "$DIST_PATH"
    
    echo " > copying resources"
    if [ ! -d "$PROJECT_PATH/android" ]; then
        cp -rp "$BUILDER_PATH/android" "$PROJECT_PATH/android"        
    fi
    cp -r "$PROJECT_PATH/android/AndroidManifest.xml" $PROJECT_PATH

    echo " > Building"
    cd $PROJECT_PATH >/dev/null
    gomobile build -tags i386 -target=android -o "$DIST_PATH/$(basename $PROJECT_PATH).apk"    
    
    if [ "$?" -eq "0" ]; then                
        echo " > Build done in $DIST_PATH"  
        rm -f "$PROJECT_PATH/AndroidManifest.xml"
    else
        echo " > Build failed" 1>&2
        rm -rf "$DIST_PATH"
        rm -f "$PROJECT_PATH/AndroidManifest.xml"
    fi

    cd - > /dev/null
}

# Build IOS
BuildIOS () {
    echo "ERROR : No implemented yet"
    # echo " > init gomobile"
    # command -v gomobile >/dev/null 2>&1 || { echo "ERROR : gomobile command not found in PATH" 1>&2; exit 1; }
    # gomobile init
    exit 3
}

# Build browser
BuildBrowser () {
    echo " > cleaning"
    rm -rf "$DIST_PATH"
    mkdir -p "$DIST_PATH"

    echo " > building"
    cd $PROJECT_PATH >/dev/null
    GOOS=js GOARCH=wasm go build -o "$DIST_PATH/main.wasm"    
    
    if [ "$?" -eq "0" ]; then
        echo " > copying resources"
        if [ ! -d "$PROJECT_PATH/browser" ]; then
            cp -rp "$BUILDER_PATH/browser" "$PROJECT_PATH"
        fi
        cp -rp "$PROJECT_PATH/browser" $(dirname $DIST_PATH)
        cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" $DIST_PATH

        echo " > Build done in $DIST_PATH"
    else
        echo " > Build failed" 1>&2
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
    echo "ERROR : Unsupported target : $TARGET" 1>&2
fi

