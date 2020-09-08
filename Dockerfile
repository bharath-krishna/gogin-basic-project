ARG GO_VERSION=1.12
FROM golang:${GO_VERSION}-alpine as builder
RUN apk --no-cache --update upgrade && apk add --no-cache ca-certificates git curl
WORKDIR /src
COPY ./ ./
RUN CGO_ENABLED=0 go build -o /app .

# Final stage: the running container.
FROM scratch AS final
COPY --from=builder /app /app
ENV http_proxy=""
ENV https_proxy=""
ENTRYPOINT ["/app"]