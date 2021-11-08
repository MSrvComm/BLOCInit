FROM golang:alpine as builder
RUN mkdir /build
WORKDIR /build
ADD . /build/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o micoinit .

FROM alpine:latest

RUN apk add iptables libcap
RUN touch /run/xtables.lock && chmod 0666 /run/xtables.lock
RUN setcap cap_net_raw,cap_net_admin+eip /sbin/xtables-legacy-multi
COPY --from=builder /build/micoinit /usr/local
ENTRYPOINT ["/usr/local/micoinit"]