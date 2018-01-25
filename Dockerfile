FROM golang:1.9-alpine

ADD ./ /go/src/github.com/micahhausler/simple-raft

WORKDIR /go/src/github.com/micahhausler/simple-raft

RUN go install

EXPOSE 3000

ENTRYPOINT ["/go/bin/simple-raft"]
