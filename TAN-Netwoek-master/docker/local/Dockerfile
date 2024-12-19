FROM golang:1.20-alpine AS builder

WORKDIR /TAN-node

ADD go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o TAN-node main.go

FROM alpine:latest AS runner

RUN apk --no-cache add ca-certificates jq

WORKDIR /TAN-node

COPY --from=builder /TAN-node/TAN-node ./
COPY ./docker/local/TAN-node.sh ./
COPY ./core-contracts/artifacts ./core-contracts/artifacts

# Expose json-rpc, libp2p and grpc ports
EXPOSE 8545 9632 1478 5001

ENTRYPOINT ["./TAN-node.sh"]
