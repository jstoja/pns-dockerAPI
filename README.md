Plug'n'Stream DockerAPI
=============

PlugAndStream Server API to launch Docker container at will.

This little server permits to create RTMP servers by launching containers through an API.

##Install the image

For the server to run, you need to create the Docker image on the machine running the API:

```
➜  ~  docker build -t pns:server .
```
Note: On Windows and OSX you need to install boot2docker. [More information here...](http://docs.docker.io/en/latest/installations)



##API

###Create a container

```
➜  ~  curl http://localhost:3001/new/stream_conference_lyon
{"Id":1,"Name":"stream_conference_lyon","PortRTMP":1935,"PortFLV":1234}%
```

###List containers
```
➜  ~  curl http://localhost:3001/list
* &{1 stream_conference_lyon 1935 1234}
* <nil>
* <nil>
* <nil>
* <nil>
```

###Delete a container
```
➜  ~  curl http://localhost:3001/delete/1
{}%
```
