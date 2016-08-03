FROM golang
ADD . /go/src/github.com/jboursiquot/packman
RUN go install github.com/jboursiquot/packman/cmd/packman
ENTRYPOINT /go/bin/packman
EXPOSE 8080
