#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd $DIR/../cmd/kl
    export KL_IAAS=azure
    go build
    ./kl digest .
popd