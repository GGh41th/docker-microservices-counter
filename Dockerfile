# Étape 1 : Construction du binaire (Builder)
FROM golang:alpine AS builder

# Installation de UPX pour compresser le binaire
RUN apk add --no-cache upx

WORKDIR /app

COPY go.mod ./
COPY main.go .

# Compilation statique sans debug + compression UPX
RUN go get github.com/redis/go-redis/v9@v9.5.1 && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o counter-app . && \
    upx --best --lzma counter-app

# Étape 2 : Image finale ultra-légère (Alpine)
FROM alpine:3.20

# Utilisateur non-root pour la sécurité
RUN adduser -D -u 1000 appuser

WORKDIR /app

# Copie unique du binaire compressé
COPY --from=builder /app/counter-app .

RUN chown appuser:appuser /app/counter-app
USER appuser

EXPOSE 5000

CMD ["./counter-app"]
