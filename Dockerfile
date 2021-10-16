FROM golang:latest AS build

#RUN go get github.com/jmunozal/goshowroom/helloworld-server

WORKDIR /go/src/github.com/jmunozal/goshowroom/helloworld-server

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

# STAGE 2: Deployment
FROM alpine

EXPOSE 8080

USER nobody:nobody
COPY --from=build /go/src/github.com/jmunozal/goshowroom/helloworld-server/server /server

CMD [ "/server" ]