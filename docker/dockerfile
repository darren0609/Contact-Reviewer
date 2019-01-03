FROM golang:latest

RUN mkdir -p /go/src/connect

WORKDIR /go/src/connect

RUN go get github.com/gorilla/mux

RUN go get github.com/google/uuid

RUN go get golang.org/x/oauth2

COPY . /go/src/connect

RUN go install

CMD /go/bin/connect

EXPOSE 8080