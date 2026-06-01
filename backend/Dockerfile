FROM golang:1.26-alpine
WORKDIR /app

ARG SECRET=replace-me
ARG APPNAME=metricraft
ARG ALLOWED_ORIGINS=http://localhost
ARG DATABASE_USERS=postgresql://postgres:password@localhost:5432/postgres?sslmode=disable
ARG GOOGLE_APP_PASSWORD=qaaz qqqv lgbm fgdg

ENV SECRET=$SECRET
ENV APPNAME=$APPNAME
ENV ALLOWED_ORIGINS=$ALLOWED_ORIGINS
ENV GOOGLE_APP_PASSWORD=$GOOGLE_APP_PASSWORD
ENV DATABASE_USERS=$DATABASE_USERS
ENV DATABASE_LOGS=postgresql://postgres:password@postgresql:5432/postgres?sslmode=disable

COPY . ./
RUN go mod download
RUN go build -o backend

EXPOSE 8080

CMD ["./backend"]
