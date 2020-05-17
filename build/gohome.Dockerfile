FROM golang:1.14.3-buster AS build

RUN mkdir -p /opt/gohome

COPY . /opt/gohome

WORKDIR /opt/gohome

RUN make test

RUN make gohome

EXPOSE 8053

ENTRYPOINT [ "/opt/gohome/gohome" ]