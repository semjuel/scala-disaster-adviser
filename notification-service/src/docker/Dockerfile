FROM alpine/git as clone
WORKDIR /app
RUN git clone https://github.com/AndriiBraslavskyi/replication.git

FROM maven:3.6.0-jdk-8-slim AS build
WORKDIR /app
COPY --from=clone /app/replication /app
RUN mvn -Dmaven.test.skip=true clean package

FROM openjdk:8u252-jre-slim
WORKDIR /app
RUN mkdir conf
COPY --from=build /app/target/replication-0.0.1-SNAPSHOT.jar /app

ARG SERVER_RETRY
ARG SERVER_TIMEOUT
ARG CLIENT_HOSTS
ARG RETRY_PERIOD
ARG NODES_FOR_QUORUM
ARG CRASHED_PING_NUMBER

EXPOSE 8080

CMD ["java", "-jar", "replication-0.0.1-SNAPSHOT.jar", "--replicas.retry-number=${SERVER_RETRY}", "--replicas.timeout=${SERVER_TIMEOUT}", "--retry-period=${RETRY_PERIOD}", "--replicas.hosts=${CLIENTS_HOSTS}", "--health.nodes-for-quorum=${NODES_FOR_QUORUM}", "--health.crashed-ping-number=${CRASHED_PING_NUMBER}"]