#!/bin/sh

# Env
if [ $# -ne 2 ] ; then
    echo "Synopsis : build.sh PROJECT_PATH TARGET"
    echo "  Available targets  : desktop android ios browser"
    exit 1
fi

if [ -z "$GOPATH" ] ; then
    echo "Error : GOPATH is not set"
    exit 1
fi

PROJECT_PATH="$1"
TARGET="$2"

if [ ! -d "$PROJECT_PATH" ] ; then
    echo "Project path not found : $PROJECT_PATH"
    exit 2
fi

DIST_PATH="$(cd "$(dirname "$PROJECT_PATH")"; pwd)/$(basename "$PROJECT_PATH")/dist"

# Common clean
rm -rf "$DIST_PATH"
mkdir "$DIST_PATH"

# Build
if [ $TARGET == "desktop" ] ; then
    echo "Building $PROJECT_PATH for desktop ..."  
    cd $PROJECT_PATH
    go build -o "$DIST_PATH/$(basename $PROJECT_PATH)"
    cd -
    echo "Build success in $DIST_PATH"  
elif [ $TARGET == "android" ] ; then
    echo "Building $PROJECT_PATH for android ..."
elif [ $TARGET == "ios" ] ; then
    echo "Building $PROJECT_PATH for ios ..."
elif [ $TARGET == "browser" ] ; then
    echo "Building $PROJECT_PATH for browser ..."
    # cd $PROJECTS_PATH/$PROJECT
    # GOARCH=js GOOS=browser go build
    # cd -
else
    echo "Unsupported target : $TARGET"
fi

