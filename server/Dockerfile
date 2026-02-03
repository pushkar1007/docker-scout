FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/scout ./cmd/scout

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/scout /app/scout
COPY web/dashboard.html /app/web/dashboard.html
ENV DOCKER_SCOUT_ADDR=:3000
EXPOSE 3000
ENTRYPOINT ["/app/scout"]
