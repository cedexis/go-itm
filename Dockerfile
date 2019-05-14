FROM golang:1.11.6-alpine3.9

RUN apk add --update \
    bash \
    bash-completion \
    build-base \
    shadow

RUN usermod --shell /bin/bash root
ADD docker/bashrc.sh /root/.bashrc

# Set the default working directory
WORKDIR /go-itm
