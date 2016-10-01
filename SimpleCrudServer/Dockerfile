FROM ubuntu:16.04

RUN apt-get update
RUN apt-get install -y golang
RUN apt-get install -y git

RUN mkdir -p /opt/lib/golang

ENV GOPATH /opt/lib/golang/
ENV PATH $GOPATH/bin:$PATH

RUN go get "github.com/gin-gonic/gin"
RUN go get "github.com/mediocregopher/radix.v2/..."

RUN mkdir -p /opt/SimpleCrudServer
COPY server.go /opt/SimpleCrudServer
WORKDIR /opt/SimpleCrudServer
RUN go build server.go
RUN rm server.go  # remove source code after app is builded

EXPOSE 3000

CMD ["./server"]