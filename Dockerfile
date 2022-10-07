FROM golang:1.19-alpine AS builder

WORKDIR /build

# Install the application dependencies.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the sources and build the application.
COPY . .
RUN CGO_ENABLED=0 go build -o /build/color-it .

#####

FROM alpine:3.16
RUN apk add ca-certificates

WORKDIR /app/

RUN addgroup --gid 1001 -S app && \
    adduser -G app --shell /bin/false --disabled-password -D -H --uid 1001 app

COPY --from=builder /build/color-it .

USER app:app

ENTRYPOINT ["/app/color-it"]
