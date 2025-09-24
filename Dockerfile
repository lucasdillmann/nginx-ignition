FROM alpine:3 AS workspace

ARG TARGETPLATFORM
WORKDIR /workspace

RUN mkdir data frontend migrations

COPY ./database/common/migrations/scripts /workspace/migrations
COPY ./frontend/build /workspace/frontend
COPY build/${TARGETPLATFORM} /workspace/nginx-ignition

FROM alpine:3

ARG NGINX_IGNITION_VERSION

EXPOSE 8090
EXPOSE 80

ENV NGINX_IGNITION_NGINX_BINARY_PATH="/usr/sbin/nginx" \
    NGINX_IGNITION_SERVER_FRONTEND_PATH="/opt/nginx-ignition/frontend" \
    NGINX_IGNITION_DATABASE_DRIVER="sqlite" \
    NGINX_IGNITION_DATABASE_MIGRATIONS_PATH="/opt/nginx-ignition/migrations" \
    NGINX_IGNITION_DATABASE_DATA_PATH="/opt/nginx-ignition/data" \
    NGINX_IGNITION_VERSION="${NGINX_IGNITION_VERSION}" \
    GOMEMLIMIT="96MiB"

ENTRYPOINT ["/opt/nginx-ignition/nginx-ignition"]
WORKDIR /opt/nginx-ignition

RUN apk update && \
    apk upgrade && \
    apk add \
      nginx \
      nginx-mod-http-js \
      nginx-mod-http-lua \
      nginx-mod-stream && \
    apk cache clean

COPY --from=workspace /workspace /opt/nginx-ignition
