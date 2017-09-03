FROM golang

LABEL maintainer="興怡"
ENV DEBIAN_FRONTEND noninteractive
ARG VNC_PASSWORD=secret
ENV VNC_PASSWORD ${VNC_PASSWORD}
ENV NAME cefgo
ENV SRC_REPO github.com/patterns
ENV PROJ_DIR /go/src/$SRC_REPO
ENV CEF_BUILD cef_binary_3.3029.1608.g27a32db_linux64
ENV OUTFILE $PROJ_DIR/cefcapi/Release/cgotest

RUN apt-get update \
 && apt-get install -y \
        supervisor x11vnc xvfb \
        dbus-x11 x11-utils \
        dwm stterm suckless-tools \
        sudo \
        build-essential libgtk2.0-dev libgtkglext1-dev \
        libglib2.0-dev libgtksourceview2.0-dev libgcrypt-dev \
        libnss3-dev libgconf-2-4 libasound2-dev \
        clang \
 && useradd --create-home --shell /bin/bash $NAME \
 && usermod --groups sudo $NAME \
 && echo "$NAME:$NAME" | chpasswd \
 && mkdir -p $PROJ_DIR \
 && curl http://opensource.spotify.com/cefbuilds/$CEF_BUILD.tar.bz2 \
        -o /go/src/cefstd.tar.bz2 \
 && cd /go/src \
 && tar xjf /go/src/cefstd.tar.bz2 \
 && mv $CEF_BUILD $PROJ_DIR/cefstd \
 && go get github.com/mattn/go-gtk/gtk \
 && rm -rf /var/lib/apt/lists/* /go/src/cefstd.tar.bz2 \
 && mkdir -p /etc/supervisor/conf.d \
 && x11vnc -storepasswd $VNC_PASSWORD /etc/vncsecret \
 && chmod 444 /etc/vncsecret 

WORKDIR $PROJ_DIR
RUN git clone https://$SRC_REPO/cefcapi \
 && mkdir cefcapi/Release \
 && cp -R cefstd/Release/* cefcapi/Release \
 && cp -R cefstd/Resources/* cefcapi/Release 

COPY . cefcapi
RUN cd cefcapi \
 && go build -ldflags "-r $PROJ_DIR/cefcapi/Release" -o $OUTFILE \
 && chown -R $NAME /go/src

COPY supervisord.conf /etc/supervisor/conf.d
CMD ["/usr/bin/supervisord","-c","/etc/supervisor/conf.d/supervisord.conf"]

EXPOSE 5900
USER $NAME
