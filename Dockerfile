FROM arm32v7/alpine

LABEL maintainer="Ismael Fernandez <fernandez.molina.ismael@gmail.com>"

WORKDIR "/go/src/github.com/ismferd/serf-publisher"

RUN apk --no-cache add tini
RUN apk --no-cache add bash

ENTRYPOINT ["/sbin/tini", "--", "./serf-publisher"]

COPY build/serf-publisher serf-publisher
