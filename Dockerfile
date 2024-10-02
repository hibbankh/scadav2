# # Dockerfile by Haziq
# FROM golang:alpine

# LABEL version = "1.0"

# LABEL maintainer="Haziq Haikal <hazique@vectolabs.com>"

# # RUN mkdir -p /var/www/image

# # COPY test.jpg /var/www/image

# WORKDIR $GOPATH/src/gitlab.com/oga

# COPY . .

# RUN go get -d -v ./...

# RUN go install -v ./...

# EXPOSE 9991

# EXPOSE 3000

# CMD [ "cmd" ]




FROM golang:1.18
LABEL version="1.0.0"
LABEL maintainer="Wan Mohd Asyraf <asyraf@vectolabs.com>"
LABEL updated_date="2023-10-26"

WORKDIR /src

# Copy application main modules
COPY go.mod go.sum ./
RUN go mod download

# Copy application modules
WORKDIR /src/app
COPY ./app/go.mod ./app/go.sum ./
RUN go mod download

# Copy application framework modules
WORKDIR /src/framework
COPY ./framework/go.mod ./framework/go.sum ./
RUN go mod download

WORKDIR /src

# Copy Source File
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /build

EXPOSE 9999
EXPOSE 3000

CMD ["/build"]