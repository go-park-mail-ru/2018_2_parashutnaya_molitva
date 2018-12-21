#!/bin/bash

tag="kekmatego"
logPath="/var/log/kekmate"

go build -o auth cmd/auth/auth.go
go build -o core cmd/server/main.go
go build -o game cmd/game/main.go

./auth ${tag} >> ${logPath}/log.auth &
./core -port 3334 $tag  >> ${logPath}/log.core &
./game -port 3336 $tag  >> ${logPath}/log.game &
