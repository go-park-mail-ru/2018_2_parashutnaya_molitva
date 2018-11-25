#!/bin/bash

if [ "$AUTH" == "yes" ]; then
    mongoClient="mongo $MONGODB_APPLICATION_DATABASE -u $MONGODB_APPLICATION_USER -p $MONGODB_APPLICATION_PASS"
else
    mongoClient="mongo $MONGODB_APPLICATION_DATABASE"
fi

echo "=> Creating a ${MONGODB_APPLICATION_DATABASE} database user with a password in MongoDB"
echo "Using $MONGODB_APPLICATION_DATABASE database"
$mongoClient << EOF
db.createCollection("dialog")
db.dialog.createIndex({"to_login":1, "from_login":1}, {unique: true})
db.createCollection("message")
db.message.createIndex({"time":1, "from_login":1, "to_login"})
db.message.createIndex({"time":1, "to_login"})
EOF

sleep 1

touch /data/db/.mongodb_init_database