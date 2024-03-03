#! /bin/bash

if [ "$1" == "w" ]; then
  env GOOS=windows GOARCH=amd64 go build -o app.exe .
else 
  go build -o app .
fi
