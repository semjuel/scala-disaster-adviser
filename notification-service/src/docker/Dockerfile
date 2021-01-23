FROM alpine/git as clone
WORKDIR /app
RUN git clone https://github.com/semjuel/scala-disaster-adviser.git

FROM maven:3.6.0-jdk-8-slim AS build
WORKDIR /app
COPY --from=clone /app/scala-disaster-adviser/notification-service /app
RUN mvn -Dmaven.test.skip=true clean package

FROM openjdk:8u252-jre-slim
WORKDIR /app
RUN mkdir conf
COPY --from=build /app/target/notification-service-1.0-SNAPSHOT.jar /app

ARG KAFKA_SERVER
ARG USER_SERVICE_ENDPOINT
ARG DISASTER_SERVICE_ENDPOINT
ARG LOCATION_GAP


EXPOSE 8080

CMD ["java", "-jar", "notification-service-1.0-SNAPSHOT.jar", "--kafka.bootstrap-servers=${KAFKA_SERVER}", "--user-service.endpoint=${USER_SERVICE_ENDPOINT}", "--disaster-service.endpoint=${DISASTER_SERVICE_ENDPOINT}",
 "--notification-service.coordinates-gap=${LOCATION_GAP}"]
