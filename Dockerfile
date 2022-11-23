FROM node:18-alpine as frontend

RUN npm install -g pnpm

WORKDIR /app
COPY client ./client/

WORKDIR /app/client
RUN pnpm install --frozen-lockfile && pnpm build


FROM --platform=linux/amd64 golang:1.19 as backend

WORKDIR /app

COPY --from=frontend /app/dist ./dist

COPY pkg/ ./pkg/
COPY go.mod go.sum Makefile main.go ./
RUN go mod download

RUN make build


FROM --platform=linux/amd64 debian:bullseye-slim as final

WORKDIR /app

COPY --from=backend /app/build/ /app/

EXPOSE 8080

CMD ["/app/server"]
