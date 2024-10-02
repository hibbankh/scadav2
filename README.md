# Project Template

## Getting started

To make it easy for you to get started with development run command below:

### Install project framework packages and modules

 ```
 ./init.sh
 ```
### Manual Setup

If there is an issue with init script run these command:

#### Clone framework and initialize new project repo

```
# clone with ssh
git clone git@gitlab.com:vecto/sw/internal/framework.git

# clone with https
git clone https://gitlab.com/vecto/sw/internal/framework.git
```

##### Initialize with windows CMD
```
del /F /Q .git
copy sample.env .env

git init
```

##### Initialize with Bash
```
rm -rf .git
cp sample.env .env

git init
```

#### Install Project and Framework dependencies

```
go get . && cd framework && go get . && cd ..
```

## Start

To start project run:
```
go run main.go
```

## Optional

### Live Reload With [Air](https://github.com/cosmtrek/air/blob/master/README.md?plain=1)

#### Installation

##### Via `go install` (Recommended)

With go 1.18 or higher:

```bash
go install github.com/cosmtrek/air@latest
```

##### Via install.sh

```bash
### binary will be $(go env GOPATH)/bin/air
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

### or install it into ./bin/
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

air -v
```

##### Via [goblin.run](https://goblin.run)

```sh
### binary will be /usr/local/bin/air
curl -sSfL https://goblin.run/github.com/cosmtrek/air | sh

### to put to a custom path
curl -sSfL https://goblin.run/github.com/cosmtrek/air | PREFIX=/tmp sh
```

##### Docker/Podman

Please pull this docker image [cosmtrek/air](https://hub.docker.com/r/cosmtrek/air).

```bash
docker/podman run -it --rm \
    -w "<PROJECT>" \
    -e "air_wd=<PROJECT>" \
    -v $(pwd):<PROJECT> \
    -p <PORT>:<APP SERVER PORT> \
    cosmtrek/air
    -c <CONF>
```
