#!/bin/bash
set -euo pipefail
project_name=$(echo $GIT_URL | awk -F/ '{print $NF}' | sed 's/.git//g')
mkdir -p $GOPATH/src/github.com/Microsoft
mkdir -p $GOPATH/src/github.com/Microsoft/${project_name}
cp -r $WORKSPACE/* $GOPATH/src/github.com/Microsoft/${project_name}/
exit_code=0

pushd $GOPATH/src/github.com/Microsoft/${project_name}
    rm -rf ./src

    ./scripts/validate_go.sh || exit_code=1
    go get -u golang.org/x/lint/golint
    go get -u github.com/onsi/ginkgo/ginkgo
    go get -u github.com/onsi/gomega/...
    export PATH="$PATH:$GOPATH/bin"
    ginkgo -race -r || exit_code=1
    exit $exit_code
popd