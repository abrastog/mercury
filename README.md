# kbridge
Bridge app for KITE Trading API. Records Live Datafeed. Publishes market tick over mangos/nanomsg socket.

# compile
```
C:\source\repo>git clone https://github.com/sivamgr/kbridge.git
C:\source\repo>cd kbridge
C:\source\repo\kbridge>go get
C:\source\repo\kbridge>go install
```

# configure and run.
modify/rename sampleconf.yml as needed.
```
C:\source\repo\kbridge>kbridge.exe -conf=c:\source\repo\kbridge\sampleconf.yml
```

# docker

Edit the sampleconf.xml with proper configuration before running the docker build

```
docker build -t kbridge .
docker run -it kbridge
```
