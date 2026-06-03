FROM golang:1.26-alpine

ARG SECRET=replace-me
ARG APPNAME=metricraft
ARG DATABASE_USERS=postgresql://postgres:password@localhost:5432/postgres?sslmode=disable
ARG GOOGLE_APP_PASSWORD=qaaz qqqv lgbm fgdg

ENV SECRET=$SECRET
ENV APPNAME=$APPNAME
ENV GOOGLE_APP_PASSWORD=$GOOGLE_APP_PASSWORD
ENV DATABASE_USERS=$DATABASE_USERS
ENV DATABASE_LOGS=postgresql://postgres:password@postgresql:5432/postgres?sslmode=disable

WORKDIR /app
COPY go.work go.work.sum ./
COPY proto ./proto
COPY backend ./backend

WORKDIR /app/backend
RUN go mod download
RUN go build -o /app/bin/backend ./cmd

EXPOSE 8080

CMD ["/app/bin/backend"]
