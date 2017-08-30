FROM ubuntu

LABEL maintainer="興怡"
ENV DEBIAN_FRONTEND noninteractive
ARG VNC_PASSWORD=secret
ENV VNC_PASSWORD ${VNC_PASSWORD}
ENV UNAME cefx
ENV PROJ_DIR /home/$UNAME
ENV CEF_BUILD cef_binary_3.3029.1608.g27a32db_linux64

RUN apt-get update \
 && apt-get install -y \
        supervisor x11vnc xvfb \
        dbus-x11 x11-utils \
        vim-tiny \
        dwm stterm suckless-tools \
        git sudo \
        build-essential libgtk2.0-dev libgtkglext1-dev \
        libnss3-dev libgconf-2-4 libasound2-dev \
        cmake ninja-build \
 && useradd --create-home --shell /bin/bash $UNAME \
 && usermod --groups sudo $UNAME \
 && echo 'cefx:cefx' | chpasswd \
 && curl http://opensource.spotify.com/cefbuilds/$CEF_BUILD.tar.bz2 \
        -o $PROJ_DIR/cefstd.tar.bz2 \
 && cd $PROJ_DIR \
 && tar xjf cefstd.tar.bz2 \
 && mv $CEF_BUILD $PROJ_DIR/cefstd \
 && rm -rf /var/lib/apt/lists/* $PROJ_DIR/cefstd.tar.bz2 \
 && mkdir -p /etc/supervisor/conf.d \
 && x11vnc -storepasswd $VNC_PASSWORD /etc/vncsecret \
 && chmod 444 /etc/vncsecret 

WORKDIR $PROJ_DIR
RUN git clone https://github.com/cztomczak/cefcapi \
 && mkdir cefcapi/Release \
 && cp -R cefstd/Release/* cefcapi/Release \
 && cp -R cefstd/Resources/* cefcapi/Release \
 && chown -R $UNAME cefstd \
 && chown -R $UNAME cefcapi \
 && cd cefcapi 
## && make

COPY supervisord.conf /etc/supervisor/conf.d
CMD ["/usr/bin/supervisord","-c","/etc/supervisor/conf.d/supervisord.conf"]

EXPOSE 5900
USER $UNAME
