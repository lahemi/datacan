#!/bin/sh

[ -z "$GOPATH" ] && echo "Your $GOPATH is not set!" && exit 1

softname="datacan"

targ="$HOME/.local/share/$softname"
src=$(pwd)

[ ! -d "$targ" ] && mkdir -p "$targ"

cp -r "$src/htmls" "$targ"
cp -r "$src/styles" "$targ"

mkdir -p "$GOPATH/src/$softname"
cp *go "$GOPATH/src/$softname"

go install "$softname"

