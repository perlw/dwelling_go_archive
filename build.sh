#!/bin/bash

BUILD_DIR="build"
PRGN_NAME="dwelling"

function build {
    rm -rf ${BUILD_DIR}
    mkdir ${BUILD_DIR}

    cd ${BUILD_DIR}
    go build ${PRGN_NAME}

    if [ $? -eq 0 ]; then
        cp -R ../resources/* .
        echo "NA"
    else
        echo "-Build failed-"
        exit $?
    fi

    cd ..

    echo "-Build complete-"
}

function run {
    if [ ! -d ${BUILD_DIR} ]; then
        echo "Must build first"
        exit -1
    fi

    cd ${BUILD_DIR}
    ./${PRGN_NAME}
}

if [ $# -ne 0 ]; then
    for i in $*
    do
        if [ $i = "build" ]; then
            build
        elif [ $i = "run" ]; then
            run
        else
            echo "Only the following command(s) are accepted: [build|run|run_local]"
        fi
    done
else
    build
fi

