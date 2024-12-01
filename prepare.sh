#!/bin/bash

init_gomobile () {
  go install golang.org/x/mobile/cmd/gomobile@latest
  go install golang.org/x/mobile/cmd/gobind@latest
  go get golang.org/x/mobile/cmd/gomobile
  go get golang.org/x/mobile/cmd/gobind
  gomobile init
}

init_gomobile
gomobile bind -v -ldflags='-s -w' -androidapi 21 -o libv2raymobile.aar
mkdir -p release/android-apilatest
cp libv2raymobile.aar release/android-apilatest/