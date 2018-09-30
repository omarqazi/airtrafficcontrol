# Start with golang base image
FROM golang:latest
MAINTAINER Omar Qazi (omar@smick.co)

# Compile latest source
ADD . /go/src/github.com/omarqazi/airtrafficcontrol
RUN go get github.com/omarqazi/airtrafficcontrol
RUN go install github.com/omarqazi/airtrafficcontrol

WORKDIR /go/src/github.com/omarqazi/airtrafficcontrol

ENTRYPOINT /go/bin/airtrafficcontrol
EXPOSE 8080
