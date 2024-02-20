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

RUN apk --no-cache add ca-certificates
COPY --from=builder /root/centaurid /usr/local/bin/centaurid

RUN addgroup --gid 1025 -S composable && adduser --uid 1025 -S composable -G composable

WORKDIR /home/composable
USER composable

# rest server
EXPOSE 1317
# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657
# grpc
EXPOSE 9090

ENTRYPOINT ["centaurid"]
CMD [ "start" ]
