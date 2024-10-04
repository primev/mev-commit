FROM alpine:latest AS cert_builder
# A unique service name.
ARG service_name
ENV SERVICE_NAME=${service_name}

RUN apk --no-cache add openssl

WORKDIR /
COPY ./p2p/gen-certificates.sh .
RUN chmod +x /gen-certificates.sh
RUN /gen-certificates.sh "/C=US/O=${SERVICE_NAME} CA" "/C=US/O=${SERVICE_NAME}"


FROM golang:1.23.2 AS mev_commit_builder
WORKDIR /
COPY ./p2p ./p2p
COPY ./contracts-abi ./contracts-abi
COPY ./x ./x
WORKDIR ./p2p
RUN CGO_ENABLED=0 GOOS=linux make build


FROM alpine:latest
# Type of node: (bootnode|bidder|provider).
ARG node_type
ENV NODE_TYPE=${node_type}

COPY --from=cert_builder /server-cert.pem /server-cert.pem
COPY --from=cert_builder /server-key.pem /server-key.pem
COPY --from=mev_commit_builder /p2p/bin/mev-commit /app/mev-commit
COPY --from=mev_commit_builder /p2p/integrationtest/config/${NODE_TYPE}.yaml /config.yaml
COPY --from=mev_commit_builder /p2p/integrationtest/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 13522 13523 13524
ENTRYPOINT ["/entrypoint.sh"]

