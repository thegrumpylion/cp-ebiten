#!/bin/zsh

export GOOS=js
export GOARCH=wasm
go build -o cpebiten.wasm github.com/jakecoffman/cpebiten/logosmash
go build -o tumble/tumble.wasm github.com/jakecoffman/cpebiten/tumble
go build -o chain/chain.wasm github.com/jakecoffman/cpebiten/chain
