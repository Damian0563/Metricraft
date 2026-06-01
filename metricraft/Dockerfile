FROM node:20-alpine
WORKDIR /app

ARG SECRET=replace-me
ARG APPNAME=metricraft
ARG DOMAIN=http://localhost

ENV SECRET=$SECRET
ENV APPNAME=$APPNAME
ENV DOMAIN=$DOMAIN
ENV NITRO_PORT=8000
ENV PORT=8000
ENV HOST=0.0.0.0
ENV NITRO_HOST=0.0.0.0

COPY . ./
RUN npm ci
RUN npm run build

EXPOSE 8000

CMD ["node", ".output/server/index.mjs"]
