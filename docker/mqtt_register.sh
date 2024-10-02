#! /usr/bin/bash

SOURCE=${BASH_SOURCE[0]}
while [ -L "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR=$(cd -P "$(dirname "$SOURCE")" >/dev/null 2>&1 && pwd)
  SOURCE=$(readlink "$SOURCE")
  [[ $SOURCE != /* ]] && SOURCE=$DIR/$SOURCE # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR=$(cd -P "$(dirname "$SOURCE")" >/dev/null 2>&1 && pwd)

echo "Create MQTT Folder"
mkdir ${DIR}/mqtt5 ${DIR}/mqtt5/config ${DIR}/mqtt5/data ${DIR}/mqtt5/log
echo "allow_anonymous false
listener 1883
listener 9001
protocol websockets
persistence false
password_file /mosquitto/config/pwfile
persistence_file mosquitto.db
persistence_location /mosquitto/data/" > ${DIR}/mqtt5/config/mosquitto.conf

touch ${DIR}/mqtt5/config/pwfile

echo ""
echo ""
echo "Run those command in docker container"
echo ""

echo "replace <asyraf> with prefered username"
echo "mosquitto_passwd -c /mosquitto/config/pwfile asyraf"
echo ""
echo "to remove user: mosquitto_passwd -D /mosquitto/config/pwfile <user-name-to-delete>"
echo ""

echo "make sure to restart container: sudo docker restart mqtt5"