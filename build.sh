#!/bin/sh
GOOS=darwin GOARCH=amd64 go build -o ./build/mosaic-tetris-mac .
GOOS=windows GOARCH=amd64 go build -o ./build/mosaic-tetris-windows . 
go build -o ./build/mosaic-tatris-linux .