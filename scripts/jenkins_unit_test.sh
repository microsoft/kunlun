#!/bin/bash
set -e
project_name=$(echo $GIT_URL | awk -F/ '{print $NF}' | sed 's/.git//g')
mkdir -p $GOPATH/src/github.com/Microsoft
mkdir -p $GOPATH/src/github.com/Microsoft/${project_name}
cp -r $WORKSPACE/* $GOPATH/src/github.com/Microsoft/${project_name}/
pushd $GOPATH/src/github.com/Microsoft/${project_name}
    rm -rf ./src
    go get -u golang.org/x/lint/golint
    go get -u github.com/onsi/ginkgo/ginkgo
    go get -u github.com/onsi/gomega/...
    export PATH="$PATH:$GOPATH/bin"
    ginkgo -race -r
popd