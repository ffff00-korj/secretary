PACKAGE_NAME = secretary
PACKAGE_PATH = ./cmd/$(PACKAGE_NAME)/main.go
BUILD_PATH = ./build

build_delete:
	rm -rf build

build: build_delete
	go build -o $(BUILD_PATH)/$(PACKAGE_NAME) $(PACKAGE_PATH)

run:
	go run $(PACKAGE_PATH)
