# Start from the latest golang base image
FROM golang:latest
LABEL maintainer="Sivakumar PS <sivamgr@gmail.com>"
RUN mkdir /kbridge
RUN mkdir /kbridge/source
WORKDIR /kbridge/source
ADD . .
RUN go get -u
RUN go build -o kbridge .
#EXPOSE 8080
CMD ["/kbridge/source/kbridge","-conf=/kbridge/source/sampleconf.yml"]