FROM alpine:3.7

LABEL maintainer="Ismael Fernandez <fernandez.molina.ismael@gmail.com>"

WORKDIR "/go/src/github.com/ismferd/serf-publisher"

RUN apk --no-cache add tini=0.16.1-r0

ENTRYPOINT ["/sbin/tini", "--", "./serf-publisher"]

COPY build/iam-role-annotator-linux-amd64 iam-role-annotator
