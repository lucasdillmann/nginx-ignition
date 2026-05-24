FROM nginx:mainline-alpine AS builder

RUN apk add --no-cache \
      gcc \
      libc-dev \
      make \
      pcre2-dev \
      zlib-dev \
      openssl-dev \
      luajit-dev \
      libmaxminddb-dev \
      libxslt-dev \
      libxml2-dev \
      linux-headers \
      curl \
      git

RUN curl -fSL https://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz -o nginx.tar.gz && \
    tar -zxvf nginx.tar.gz

RUN git clone --depth 1 https://github.com/vision5/ngx_devel_kit.git /ngx_devel_kit && \
    git clone --depth 1 https://github.com/openresty/lua-nginx-module.git /lua-nginx-module && \
    git clone --depth 1 https://github.com/vozlt/nginx-module-vts.git /nginx-module-vts && \
    git clone --depth 1 https://github.com/leev/ngx_http_geoip2_module.git /ngx_http_geoip2_module && \
    git clone --depth 1 https://github.com/nginx/njs.git /njs && \
    git clone --depth 1 https://github.com/openresty/lua-resty-core.git /lua-resty-core && \
    git clone --depth 1 https://github.com/openresty/lua-resty-lrucache.git /lua-resty-lrucache

ENV LUAJIT_LIB=/usr/lib \
    LUAJIT_INC=/usr/include/luajit-2.1

RUN cd nginx-${NGINX_VERSION} && \
    ./configure \
      --with-compat \
      --with-http_ssl_module \
      --with-stream \
      --add-dynamic-module=/ngx_devel_kit \
      --add-dynamic-module=/lua-nginx-module \
      --add-dynamic-module=/nginx-module-vts \
      --add-dynamic-module=/ngx_http_geoip2_module \
      --add-dynamic-module=/njs/nginx && \
    make modules

RUN mkdir -p /modules && \
    cp nginx-${NGINX_VERSION}/objs/*.so /modules/

RUN mkdir -p /modules-lua && \
    cd /lua-resty-core && make install LUA_LIB_DIR=/modules-lua && \
    cd /lua-resty-lrucache && make install LUA_LIB_DIR=/modules-lua

FROM alpine:edge AS workspace

RUN apk add --no-cache \
      luajit \
      libmaxminddb \
      libxml2 \
      libxslt \
      pcre2 \
      zlib \
      openssl \
      ca-certificates \
      curl && \
    update-ca-certificates

ARG TARGETPLATFORM
WORKDIR /

RUN mkdir -p opt/nginx-ignition/data \
             opt/nginx-ignition/frontend \
             opt/nginx-ignition/migrations \
             usr/lib/nginx/modules \
             usr/share/lua/5.1

COPY ./database/common/migrations/scripts /opt/nginx-ignition/migrations
COPY ./frontend/build /opt/nginx-ignition/frontend
COPY build/${TARGETPLATFORM} /opt/nginx-ignition/nginx-ignition

COPY --from=builder /usr/sbin/nginx /usr/sbin/nginx
COPY --from=builder /etc/nginx/ /etc/nginx/
COPY --from=builder /modules/ /usr/lib/nginx/modules/
COPY --from=builder /modules-lua/ /usr/share/lua/5.1/

FROM scratch

ENV NGINX_IGNITION_NGINX_BINARY_PATH="/usr/sbin/nginx" \
    NGINX_IGNITION_SERVER_FRONTEND_PATH="/opt/nginx-ignition/frontend" \
    NGINX_IGNITION_DATABASE_DRIVER="sqlite" \
    NGINX_IGNITION_DATABASE_MIGRATIONS_PATH="/opt/nginx-ignition/migrations" \
    NGINX_IGNITION_DATABASE_DATA_PATH="/opt/nginx-ignition/data" \
    GOMEMLIMIT="128MiB"

EXPOSE 8090
ENTRYPOINT ["/opt/nginx-ignition/nginx-ignition"]
WORKDIR /opt/nginx-ignition
VOLUME /opt/nginx-ignition/data
HEALTHCHECK \
    --interval=5s \
    --timeout=5s \
    --retries=3 \
    CMD curl -f http://localhost:8090/api/health/liveness || exit 1

COPY --from=workspace / /
