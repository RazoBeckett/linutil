BINARY_NAME = linutil
MAIN_PACKAGE_PATH = ${pwd}

build:
	go build -o bin/${BINARY_NAME} $(MAIN_PACKAGE_PATH)

run:
	go run $(MAIN_PACKAGE_PATH)

clean:
	rm -rf bin

.PHONY: build run clean
