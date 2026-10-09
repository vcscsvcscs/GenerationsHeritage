FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/backup .

# distroless/static ships CA certificates and tzdata (needed for TZ)
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/backup /backup
USER nonroot:nonroot
ENTRYPOINT ["/backup"]
