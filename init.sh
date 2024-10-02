#! /usr/bin/bash

# Download/Replace core-framework Repo from https://gitlab.com/vecto/sw/internal/core-framework.git

CORE_FOLDER_NAME="framework"
CORE_REPO="https://gitlab.com/vecto/sw/internal/${CORE_FOLDER_NAME}.git"

SOURCE=${BASH_SOURCE[0]}
while [ -L "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR=$(cd -P "$(dirname "$SOURCE")" >/dev/null 2>&1 && pwd)
  SOURCE=$(readlink "$SOURCE")
  [[ $SOURCE != /* ]] && SOURCE=$DIR/$SOURCE # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR=$(cd -P "$(dirname "$SOURCE")" >/dev/null 2>&1 && pwd)

if [ ! -d "${DIR}/${CORE_FOLDER_NAME}" ]; then
  # Clone core-framework
  git clone $CORE_REPO
else
  # Pull core-framework
  cd $CORE_FOLDER_NAME
  git fetch && git pull
  cd $DIR
fi

rm -rf .git
cp sample.env .env

go mod download

cd "${DIR}/framework"
go mod download

cd "${DIR}/app"
go mod download

cd $DIR