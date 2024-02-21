ARG ARCH="amd64"
ARG OS="linux"

FROM golang:1.21.3

WORKDIR /ttm

COPY .build/numa_exporter /usr/bin/numa_exporter

RUN adduser numa_exporter
RUN chmod -R 755 "$GOPATH"
RUN chown -R numa_exporter:numa_exporter /numa_exporter

USER       numa_exporter

EXPOSE     10000

CMD ["/usr/bin/numa_exporter"]
