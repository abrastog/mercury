# Start from the latest golang base image
# Edit the sampleconf.yml and save as appconfig.yml before running the docker build
FROM golang:latest
LABEL maintainer="Sivakumar PS <sivamgr@gmail.com>"
RUN mkdir -p /solar/src/mercury
WORKDIR /solar/src/mercury
ADD . .
RUN go get -u
RUN go build -o mercury .
#EXPOSE 80
#EXPOSE 443
CMD ["/solar/src/mercury/mercury","-conf=/solar/src/mercury/appconfig.yml"]