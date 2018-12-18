#!/bin/bash
ssh-keyscan -H $KSENOBAIT_VSCALE_SERVER >> ~/.ssh/known_hosts
chmod 600 ./ssh_key
mkdir ./kekmate # создаем папку которая поедет на сервер

ssh -i ./ssh_key -tt travis@$KSENOBAIT_VSCALE_SERVER << EOF
cd ~/go/src/github.com/go-park-mail-ru/2018_2_parashutnaya_molitva
git fetch
git reset --hard $TRAVIS_COMMIT
git clean -f
./stop.sh
./run.sh
EOF