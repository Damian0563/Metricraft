FROM golang:1.26-alpine
WORKDIR /app

ARG SECRET
ARG APPNAME
ARG ALLOWED_ORIGINS
ARG DATABASE_USERS

ENV SECRET=$SECRET
ENV APPNAME=$APPNAME
ENV ALLOWED_ORIGINS=$ALLOWED_ORIGINS
ENV DATABASE_USERS=$DATABASE_USERS

# DATABASE_LOGS is NOT a build arg: it is injected at runtime by
# docker-compose and points at the internal "postgresql" service, so it
# never leaves the Docker network.

# Explicit deployment mode (overrides any value in a copied .env)
ENV MODE=standalone

COPY . ./
RUN go mod download
RUN go build -o backend

EXPOSE 8080

CMD ["./backend"]
