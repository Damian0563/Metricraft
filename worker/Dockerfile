FROM golang:1.26-alpine

ARG SECRET=replace-me
ARG APPNAME=metricraft

ENV SECRET=$SECRET
ENV APPNAME=$APPNAME
ENV DATABASE_LOGS=postgresql://postgres:password@postgresql:5432/postgres?sslmode=disable
ENV DEST_PORT=3000

WORKDIR /app
COPY go.work go.work.sum ./
COPY proto ./proto
COPY worker ./worker

WORKDIR /app/worker
RUN go mod download
RUN go build -o /app/bin/worker ./cmd

EXPOSE 8081
EXPOSE 50051

CMD ["/app/bin/worker"]
