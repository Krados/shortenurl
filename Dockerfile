FROM golang:1.19-alpine AS build_base

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_base AS server_builder

COPY . .

RUN go build -o ./bin/ ./...

FROM alpine

COPY --from=server_builder /app/bin/cmd /bin/cmd
COPY --from=server_builder /app/configs /configs

EXPOSE 8080

CMD [ "bin/cmd", "-conf", "./configs" ]