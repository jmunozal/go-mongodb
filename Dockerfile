FROM golang:latest AS build

RUN go get github.com/jmunozal/go-mongodb

WORKDIR /go/src/github.com/jmunozal/go-mongodb

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build . -o server

# STAGE 2: Deployment
FROM alpine

EXPOSE 8080

COPY --from=build /go/src/github.com/jmunozal/go-mongodb/server /server

CMD [ "/server" ]