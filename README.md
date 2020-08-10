# mercury
Bridge app for KITE Trading API. Records Live Datafeed. Publishes market tick over mangos/nanomsg socket.

# compile
```
C:\source\repo>git clone https://github.com/sivamgr/mercury.git
C:\source\repo>cd mercury
C:\source\repo\mercury>go get
C:\source\repo\mercury>go install
```

# configure and run.
modify/rename sampleconf.yml as needed.
```
C:\source\repo\mercury>mercury.exe -conf=c:\source\repo\mercury\appconfig.yml
```

# docker

Edit the sampleconf.xml with proper configuration before running the docker build

```
docker build -t mercury .
docker run -it mercury
```
