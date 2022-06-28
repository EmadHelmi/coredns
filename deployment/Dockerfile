FROM coredns_base:latest
LABEL maintainer="Emad Helmi"

WORKDIR /code

ADD custom_plugins plugin
ADD ../plugin.cfg .

RUN go generate
RUN go build
