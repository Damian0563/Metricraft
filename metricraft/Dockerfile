from node:20-alpine
workdir /app
copy . ./
run npm ci
run npm run build
EXPOSE 8000
#dev setup
CMD ["npm", "run", "dev"]
