# 第一阶段：构建前端
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# 使用淘宝 npm 镜像
RUN npm config set registry https://registry.npmmirror.com

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ .
RUN npm run build

# 第二阶段：构建后端
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# 使用国内 Go 模块镜像
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o myplants main.go

# 第三阶段：运行环境（基于 Alpine，兼容 Gentoo Linux）
FROM alpine:3.18

WORKDIR /app

RUN mkdir -p /app/data /app/uploads /app/logs /app/frontend/dist

COPY --from=backend-builder /app/myplants /app/myplants
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist
COPY config.yaml /app/config.yaml

RUN chmod +x /app/myplants

EXPOSE 8020

VOLUME ["/app/data", "/app/uploads", "/app/logs"]

CMD ["./myplants"]