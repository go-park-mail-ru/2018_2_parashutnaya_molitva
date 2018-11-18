#!/bin/bash

if [ "$AUTH" == "yes" ]; then
    mongoClient="mongo $MONGODB_APPLICATION_DATABASE -u $MONGODB_APPLICATION_USER -p $MONGODB_APPLICATION_PASS"
else
    mongoClient="mongo $MONGODB_APPLICATION_DATABASE"
fi

echo "=> Creating a ${MONGODB_APPLICATION_DATABASE} database user with a password in MongoDB"
echo "Using $MONGODB_APPLICATION_DATABASE database"
$mongoClient << EOF
db.createCollection("users")
db.users.createIndex({"email":1}, {unique: true})
db.fruit.createIndex({type: 1},
                      {collation: { locale: 'en', strength: 1 },
                       unique: true
                      })
EOF

sleep 1

touch /data/db/.mongodb_init_database
