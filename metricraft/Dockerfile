FROM node:20-alpine
WORKDIR /app

ARG SECRET
ARG APPNAME
ARG DOMAIN

ENV SECRET=$SECRET
ENV APPNAME=$APPNAME
# Public backend base URL the browser uses to reach the API/WebSocket,
# e.g. "https://metrics.example.com" or "http://203.0.113.10:8080".
ENV DOMAIN=$DOMAIN

# Explicit deployment mode (overrides any value in a copied .env)
ENV MODE=standalone

# Nitro (production server) reads NITRO_PORT/PORT; default would be 3000
ENV NITRO_PORT=8000
ENV PORT=8000

COPY . ./
RUN npm ci
RUN npm run build

EXPOSE 8000

CMD ["node", ".output/server/index.mjs"]
