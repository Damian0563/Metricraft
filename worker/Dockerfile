from golang:1.26-alpine
workdir /app
copy . ./
run go mod download
run go build -o worker
EXPOSE 8081
CMD ["./worker"]
