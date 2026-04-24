FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o mocker ./cmd/server

FROM nginx:alpine

COPY --from=builder /build/mocker /mocker
COPY nginx.conf /etc/nginx/nginx.conf
COPY frontend/dist /usr/share/nginx/html/

RUN apk add --no-cache bash

RUN echo '#!/bin/bash' > /start.sh && \
    echo '/mocker &' >> /start.sh && \
    echo 'nginx -g "daemon off;"' >> /start.sh && \
    chmod +x /start.sh

EXPOSE 80

CMD ["/start.sh"]