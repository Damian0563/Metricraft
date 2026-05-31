FROM golang:1.26-alpine
WORKDIR /app

ARG SECRET
ARG APPNAME

ENV SECRET=$SECRET
ENV APPNAME=$APPNAME

# DATABASE_LOGS is NOT set here on purpose: it is injected at runtime by
# docker-compose and points at the internal "postgresql" service, so it
# never leaves the Docker network.

# Explicit deployment mode (overrides any value in a copied .env)
ENV MODE=standalone

COPY . ./
RUN go mod download
RUN go build -o worker

EXPOSE 8081
EXPOSE 50051

CMD ["./worker"]
