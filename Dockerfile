FROM golang:alpine
RUN go version

ADD . /go/src/app
WORKDIR /go/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Expose 8080
# Gin will use the PORT env var
ENV PORT 8080
EXPOSE 8080

# Install git
RUN apk add --no-cache git
# Fetch deps
RUN go get
# Remove git
RUN apk del git

# Compile app
RUN go build -o main .
# Run app
CMD ["/go/src/app/main"]