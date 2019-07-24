FROM golang:alpine AS builder
RUN apk update && apk add git 
# bash build-base curl

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/meli-proxy

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
# RUN go install -v ./...

WORKDIR $GOPATH/src/meli-proxy/cmd/proxy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/proxy .

RUN chmod +x /go/bin/proxy

FROM scratch

# This container exposes port 8081 to the outside world
EXPOSE 8081

# Copy our static executable.
COPY --from=builder /go/bin/proxy /go/bin/proxy

# Run the hello binary.
ENTRYPOINT ["/go/bin/proxy"]