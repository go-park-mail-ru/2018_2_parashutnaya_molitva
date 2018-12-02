#!/bin/bash

tag="kekmatego"

go build -race -o auth cmd/auth/auth.go
go build -race -o core cmd/server/main.go
go build -race -o game cmd/game/main.go

./auth $tag > log.auth &
./core -port 3334 $tag  > log.core &
./game -port 3336 $tag  > log.game &
