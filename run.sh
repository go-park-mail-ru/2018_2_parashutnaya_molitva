#!/bin/bash

tag="kekmatego"
logPath="/var/log/kekmate"

go build -race -o auth cmd/auth/auth.go
go build -race -o core cmd/server/main.go
go build -race -o game cmd/game/main.go

./auth ${tag} >> ${logPath}/log.auth &
./core -port 3334 $tag  >> ${logPath}/log.core &
./game -port 3336 $tag  >> ${logPath}/log.game &
