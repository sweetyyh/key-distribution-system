# ---- Stage 1: Build Vue frontend ----
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend
COPY stitch_elite_wholesale_card_ui/vue_app/package.json ./
RUN npm install
COPY stitch_elite_wholesale_card_ui/vue_app/ ./
RUN npm run build

# ---- Stage 2: Build Go binary ----
FROM golang:1.23-alpine AS go-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /kds ./cmd/server

# ---- Stage 3: Final runtime image ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=go-builder /kds ./kds

# Copy static frontend assets
COPY --from=frontend-builder /app/frontend/dist ./stitch_elite_wholesale_card_ui/vue_app/dist
COPY stitch_elite_wholesale_card_ui/app ./stitch_elite_wholesale_card_ui/app
COPY stitch_elite_wholesale_card_ui/index.html ./stitch_elite_wholesale_card_ui/index.html

# Copy default config (will be overridden by volume mount in production)
COPY configs/config.yaml ./configs/config.yaml

# Create data directory for SQLite
RUN mkdir -p ./data

EXPOSE 9000

ENV KDS_CONFIG=/app/configs/config.yaml

CMD ["./kds"]
