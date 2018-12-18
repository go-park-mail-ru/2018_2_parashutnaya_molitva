#!/bin/bash
curl -L https://git.io/vp6lP | sh
export PATH=$PATH:$TRAVIS_BUILD_DIR/bin
gometalinter --config=gometalinter.json  ./internal/...