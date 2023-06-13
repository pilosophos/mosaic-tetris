#!/bin/sh
GOOS=darwin GOARCH=amd64 go build -o ./build/mosaic-tetris-mac .
GOOS=windows GOARCH=amd64 go build -o ./build/mosaic-tetris-windows.exe . 
go build -o ./build/mosaic-tetris-linux .