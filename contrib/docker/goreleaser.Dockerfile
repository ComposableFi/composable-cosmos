FROM golang:1.20 AS builder

WORKDIR /root
COPY ./dist/ /root/

ARG TARGETARCH
RUN if [ "${TARGETARCH}" = "arm64" ]; then \
  cp linux_linux_arm64/centaurid /root/centaurid; \
  else \
  cp linux_linux_amd64_v1/centaurid /root/centaurid; \
  fi

FROM alpine:latest

WORKDIR /root
RUN apk --no-cache add ca-certificates
COPY --from=builder /root/centaurid /usr/local/bin/centaurid

ENTRYPOINT ["centaurid"]
CMD [ "start" ]
