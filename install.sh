#!/bin/bash

realpath() {
    [[ $1 = /* ]] && echo "$1" || echo "$PWD/${1#./}"
}

REALPATH=`realpath $0`
REALPATH=`dirname $REALPATH`

if [ -z $GOPATH ]; then
    mkdir -p '.gostuff'
    export GOPATH="$REALPATH/.gostuff"
    echo "No GOPATH, setting a default one"
fi

echo "GOPATH=$GOPATH"


go get -v -x github.com/jmckaskill/gospdy

if [ $? != 0 ]; then
    echo
    echo 'You need the `go` command line tool'
    echo "More information at http://www.golang.org"
    echo

    # ubuntu: apt-get install golang golang-go

    # mac os: http://www.golang.org/

    exit 1
fi

echo
echo 'You can now use `make` to rebuild'
echo

make

