FROM golang:1.23 AS dev
WORKDIR /app
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/air-verse/air@latest
RUN go install go.uber.org/mock/mockgen@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN export PATH=$PATH:$(go env GOPATH)/bin
COPY ./ ./

FROM golang:1.23 AS builder
WORKDIR /app
COPY --from=dev /app /app
RUN go mod download
ENV PORT=8080
RUN CGO_ENABLED=0 go build -o /gin_server -ldflags="-s -w" -trimpath


FROM gcr.io/distroless/static-debian12 AS production
WORKDIR /app
COPY --from=builder /app /app
COPY --from=builder /gin_server /
ENV PORT=8080
ENTRYPOINT [ "/gin_server" ]
