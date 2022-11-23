.PHONY: clean build

APP_NAME = server
BUILD_DIR = $(PWD)/build

clean:
	rm -rf ./build
	rm -rf ./dist

build:
	# build the backend
	go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

build-frontend:
	# build the frontend
	cd client && pnpm build
