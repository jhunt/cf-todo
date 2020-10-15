FROM golang:1.15-alpine3.12 AS go

WORKDIR /go/src/cf-todo
COPY . .

RUN go build .

FROM node:14.5.0 AS node
WORKDIR /app
COPY . .

RUN cd ux/ \
 && yarn add @vue/cli \
 && yarn install \
 && yarn build

FROM alpine:3.12
COPY --from=go    /go/src/cf-todo/cf-todo /usr/bin/cf-todo
COPY --from=node  /app/ux/dist /var/lib/webroot

ENV BIND=:8080
ENV WEBROOT=/var/lib/webroot
CMD ["cf-todo"]
