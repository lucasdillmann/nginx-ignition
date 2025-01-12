FROM golang:1.23 AS builder

WORKDIR /nginx-ignition
COPY . /nginx-ignition

ENV CGO_ENABLED=0

RUN go work sync && cd application && go build

FROM alpine:3

EXPOSE 8090:8090
EXPOSE 80:80

ENV NGINX_IGNITION_NGINX_BINARY_PATH="/usr/sbin/nginx"
ENV NGINX_IGNITION_FRONTEND_PATH="/opt/nginx-ignition/frontend"
ENV NGINX_IGNITION_DATABASE_DRIVER="sqlite"
ENV NGINX_IGNITION_DATABASE_MIGRATIONS_PATH="/opt/nginx-ignition/migrations"
ENV NGINX_IGNITION_DATABASE_DATA_PATH="/opt/nginx-ignition/data"

ENTRYPOINT ["/opt/nginx-ignition/nginx-ignition"]
WORKDIR /opt/nginx-ignition

RUN mkdir data frontend migrations && \
    apk update && \
    apk add nginx nginx-mod-http-js nginx-mod-http-lua && \
    apk cache clean

COPY ./database/common/migrations/scripts /opt/nginx-ignition/migrations
COPY ./frontend/build /opt/nginx-ignition/frontend
COPY --from=builder /nginx-ignition/application/application /opt/nginx-ignition/nginx-ignition
