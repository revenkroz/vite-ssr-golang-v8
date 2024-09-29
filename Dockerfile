FROM node:22-alpine AS frontend

WORKDIR /app/frontend

RUN --mount=type=bind,source=frontend/package.json,target=package.json \
    --mount=type=bind,source=frontend/yarn.lock,target=yarn.lock \
    --mount=type=cache,target=/root/.yarn \
    yarn install --frozen-lockfile

COPY frontend .

RUN yarn build


FROM --platform=linux/amd64 golang:1.22 AS backend

RUN apt-get update && apt-get install -y \
    gcc \
    g++

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY main.go .
COPY ./pkg ./pkg

COPY --from=frontend /app/dist ./dist

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
RUN go build -o build/server -ldflags "-w -s" .


FROM --platform=linux/amd64 debian:12 AS final

WORKDIR /app

COPY --from=backend /app/build/ /app/

EXPOSE 8080

CMD ["/app/server"]
