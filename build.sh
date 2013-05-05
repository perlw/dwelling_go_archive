#!/bin/bash

BUILD_DIR="build"
PRGN_NAME="dwelling"

function build {
    echo "-Building ${PRGN_NAME}-"

    echo "Cleaning up"
    rm -rf ${BUILD_DIR}
    mkdir ${BUILD_DIR}

    echo "Building..."
    cd ${BUILD_DIR}
    go build ${PRGN_NAME}

    if [ $? -eq 0 ]; then
        echo "Build OK"

        echo "Running tests..."
        tests
        if [ $? -eq 0 ]; then
            echo "OK"
        else
            echo "FAIL"
            exit $?
        fi

        echo "Copying resources"
        cp -R ../resources/* .
    else
        echo "Tests FAIL"
        exit $?
    fi

    cd ..

    echo "-Build complete-"
}

function tests {
    go test dwelling/matrix
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
            echo "Only the following command(s) are accepted: [build|run]"
        fi
    done
else
    build
fi

