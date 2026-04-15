from golang:1.26-alpine
workdir /app
copy . ./
run go mod download
run go build -o backend
EXPOSE 8080
CMD ["./backend"]
