FROM golang:1.14.3-buster AS build

RUN mkdir -p /opt/gohome

COPY . /opt/gohome

WORKDIR /opt/gohome

RUN make test

RUN make gohome

FROM debian:buster-slim AS run

RUN mkdir -p /opt

COPY --from=build /opt/gohome/gohome /opt/gohome

EXPOSE 8053

ENTRYPOINT [ "/opt/gohome" ]