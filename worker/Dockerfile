FROM golang:1.26-alpine
WORKDIR /app

ARG SECRET=replace-me
ARG APPNAME=metricraft

ENV SECRET=$SECRET
ENV APPNAME=$APPNAME
ENV DATABASE_LOGS=postgresql://postgres:password@postgresql:5432/postgres?sslmode=disable
ENV DEST_PORT=3000

COPY . ./
RUN go mod download
RUN go build -o worker

EXPOSE 8081
EXPOSE 50051

CMD ["./worker"]
