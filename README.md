# docker-cefcapi
This is a Docker image for CEF C API experiments.
 CEF is the Chromium Embedded Framework

> [CEF](https://bitbucket.org/chromiumembedded/cef/wiki/UsingTheCAPI.md)

## Quickstart
1. Install [Docker](https://docker.com/)
2. Git clone the [src](https://github.com/patterns/docker-cefcapi/)
3. From the cloned directory, build

You can assign the VNC password for the image:

```console
$ docker build --build-arg VNC_PASSWORD=abc123 -t cefcapi .
$ docker run -it --rm -p 8088:5900 cefcapi
```

- Then launch your VNC viewer at vnc://127.0.0.1:8088
- Type the VNC_PASSWORD when prompted.
- From inside the container, it's shift-alt-enter to launch a terminal.

```console
 cd ~/cefcapi && make
```

## Credits

CEFCAPI project is by
 [Czarek Tomczak](https://github.com/cztomczak/cefcapi)

