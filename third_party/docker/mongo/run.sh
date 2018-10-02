#!/bin/bash
set -m

mongodb_cmd="mongod --bind_ip_all"
cmd="$mongodb_cmd"

if [ "$AUTH" == "yes" ]; then
    cmd="$cmd --auth"
fi

$cmd &

if [ ! -f /data/db/.mongodb_password_set ]; then
    /set_mongodb_password.sh
fi

if [ ! -f /data/db/.mongodb_init_database ]; then
    /init_database_$MONGODB_APPLICATION_DATABASE.sh
fi

fg