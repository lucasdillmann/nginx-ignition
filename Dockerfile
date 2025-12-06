FROM alpine:3 AS workspace

ARG TARGETPLATFORM
WORKDIR /workspace

RUN mkdir data frontend migrations

COPY ./database/common/migrations/scripts /workspace/migrations
COPY ./frontend/build /workspace/frontend
COPY build/${TARGETPLATFORM} /workspace/nginx-ignition

FROM alpine:3

EXPOSE 8090

HEALTHCHECK \
    --interval=5s \
    --timeout=5s \
    --retries=3 \
    CMD curl -f http://localhost:8090/api/health/liveness || exit 1

ENV NGINX_IGNITION_NGINX_BINARY_PATH="/usr/sbin/nginx" \
    NGINX_IGNITION_SERVER_FRONTEND_PATH="/opt/nginx-ignition/frontend" \
    NGINX_IGNITION_DATABASE_DRIVER="sqlite" \
    NGINX_IGNITION_DATABASE_MIGRATIONS_PATH="/opt/nginx-ignition/migrations" \
    NGINX_IGNITION_DATABASE_DATA_PATH="/opt/nginx-ignition/data" \
    GOMEMLIMIT="96MiB"

ENTRYPOINT ["/opt/nginx-ignition/nginx-ignition"]
WORKDIR /opt/nginx-ignition

RUN apk update && \
    apk upgrade && \
    apk add \
      nginx \
      nginx-mod-http-js \
      nginx-mod-http-lua \
      nginx-mod-stream \
      ca-certificates \
      curl && \
    apk cache clean && \
    update-ca-certificates

COPY --from=workspace /workspace /opt/nginx-ignition
