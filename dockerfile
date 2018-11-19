FROM golang:latest

RUN mkdir -p /go/src/connect

WORKDIR /go/src/connect

COPY . /go/src/connect

RUN go install connect

CMD /go/bin/connect

EXPOSE 8080