FROM golang:alpine AS build

WORKDIR /app

COPY . .

RUN go get ./...

RUN go build -o db-adapter

RUN apk update && apk add ca-certificates && update-ca-certificates

FROM busybox:1.36.1

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /app/db-adapter /app/

CMD [ "/app/db-adapter" ]