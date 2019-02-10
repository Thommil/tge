#!/bin/sh

PROJECTS_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [ $# -ne 2 ] ; then
    echo "Synopsis : build.sh PROJECT TARGET"
    echo "  Available projects : $(ls $PROJECTS_PATH | tr '\n' ' ' | sed 's/build.sh//')"
    echo "  Available targets  : desktop android ios browser"
    exit 1
fi

PROJECT="$1"
TARGET="$2"

if [ ! -d "$PROJECTS_PATH/$PROJECT" ] ; then
    echo "Project not found : $PROJECT"
    exit 2
fi

if [ $TARGET == "desktop" ] ; then
    echo "Building $PROJECT for desktop ..."
    cd $PROJECTS_PATH/$PROJECT
    export GOOS=darwin
    go build
    cd -
elif [ $TARGET == "android" ] ; then
    echo "Building $PROJECT for android ..."
elif [ $TARGET == "ios" ] ; then
    echo "Building $PROJECT for ios ..."
elif [ $TARGET == "browser" ] ; then
    echo "Building $PROJECT for browser ..."
else
    echo "Unsupported target : $TARGET"
fi

