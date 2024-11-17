FROM gradle:8.10.2 AS build
RUN mkdir -p /home/gradle/nginx-ignition
COPY --chown=gradle:gradle . /home/gradle/nginx-ignition
WORKDIR /home/gradle/nginx-ignition
RUN gradle assemble --no-daemon

FROM eclipse-temurin:21-jre-alpine AS runtime
EXPOSE 8090:8090
ENV NGINX_IGNITION_NGINX_BINARY_PATH="/usr/sbin/nginx"
ENV NGINX_IGNITION_DATABASE_URL="jdbc:h2:/opt/nginx-ignition/data/nginx-ignition;DB_CLOSE_DELAY=-1"
ENV NGINX_IGNITION_DATABASE_USERNAME="sa"
ENV NGINX_IGNITION_DATABASE_PASSWORD=""
ENTRYPOINT ["java", "-jar", "/opt/nginx-ignition/nginx-ignition.jar"]
RUN mkdir -p /opt/nginx-ignition/data && \
    apk update && \
    apk add nginx && \
    apk cache clean
COPY --from=build /home/gradle/nginx-ignition/application/build/libs/application-all.jar /opt/nginx-ignition/nginx-ignition.jar
