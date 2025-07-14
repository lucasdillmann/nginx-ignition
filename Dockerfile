FROM alpine:3

ARG TARGETPLATFORM
ARG NGINX_IGNITION_VERSION

EXPOSE 8090:8090
EXPOSE 80:80

ENV NGINX_IGNITION_NGINX_BINARY_PATH="/usr/sbin/nginx"
ENV NGINX_IGNITION_SERVER_FRONTEND_PATH="/opt/nginx-ignition/frontend"
ENV NGINX_IGNITION_DATABASE_DRIVER="sqlite"
ENV NGINX_IGNITION_DATABASE_MIGRATIONS_PATH="/opt/nginx-ignition/migrations"
ENV NGINX_IGNITION_DATABASE_DATA_PATH="/opt/nginx-ignition/data"
ENV NGINX_IGNITION_VERSION="${NGINX_IGNITION_VERSION}"
ENV GOMEMLIMIT="96MiB"

ENTRYPOINT ["/opt/nginx-ignition/entrypoint.sh"]
WORKDIR /opt/nginx-ignition

RUN mkdir data frontend migrations && \
    apk update && \
    apk add nginx nginx-mod-http-js nginx-mod-http-lua nginx-mod-stream tailscale && \
    apk cache clean

COPY ./docker-entrypoint.sh /opt/nginx-ignition/entrypoint.sh
RUN chmod +x /opt/nginx-ignition/entrypoint.sh

COPY ./database/common/migrations/scripts /opt/nginx-ignition/migrations
COPY ./frontend/build /opt/nginx-ignition/frontend
COPY build/${TARGETPLATFORM} /opt/nginx-ignition/nginx-ignition
