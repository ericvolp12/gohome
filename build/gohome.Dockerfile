FROM golang:1.14.3-buster AS build

RUN mkdir -p /opt/gohome

COPY . /opt/gohome

WORKDIR /opt/gohome

RUN make test

RUN make gohome-static

FROM scratch AS run

COPY --from=build /opt/gohome/gohome-static /gohome-static

EXPOSE 8053

ENTRYPOINT [ "/gohome-static" ]