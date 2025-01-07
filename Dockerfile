FROM golang:1.23 AS builder
WORKDIR /src
COPY . /src
ENV CGO_ENABLED=0
RUN go work sync && cd application && go build

FROM alpine:3
EXPOSE 8090:8090
EXPOSE 80:80
ENV NGINX_IGNITION_NGINX_BINARY_PATH="/usr/sbin/nginx"
ENV NGINX_IGNITION_FRONTEND_PATH="/opt/nginx-ignition/frontend"
ENTRYPOINT ["/opt/nginx-ignition/nginx-ignition"]
WORKDIR /opt/nginx-ignition
RUN mkdir data frontend && \
    apk update && \
    apk add nginx nginx-mod-http-js nginx-mod-http-lua && \
    apk cache clean
COPY ./frontend/build /opt/nginx-ignition/frontend
COPY --from=builder /src/application/application /opt/nginx-ignition/nginx-ignition
