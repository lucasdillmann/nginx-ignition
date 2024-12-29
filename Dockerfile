FROM eclipse-temurin:21-jre-alpine
EXPOSE 8090:8090
ENV NGINX_IGNITION_NGINX_BINARY_PATH="/usr/sbin/nginx"
ENV NGINX_IGNITION_DATABASE_URL="jdbc:h2:/opt/nginx-ignition/data/nginx-ignition;DB_CLOSE_DELAY=-1"
ENV NGINX_IGNITION_DATABASE_USERNAME="sa"
ENTRYPOINT ["java", "-Xms32m", "-Xmx128m", "-jar", "/opt/nginx-ignition/nginx-ignition.jar"]
RUN mkdir -p /opt/nginx-ignition/data && \
    apk update && \
    apk add nginx nginx-mod-http-js nginx-mod-http-lua && \
    apk cache clean
COPY application/build/libs/application-0.0.0-all.jar /opt/nginx-ignition/nginx-ignition.jar
