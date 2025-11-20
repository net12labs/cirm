# filepath: /home/lxk/Desktop/cirm/bins/vf_test/Makefile
.PHONY: all build run clean test

all: build run

build:
	@echo "Building test application..."
	go build -o vf_test .

run: build
	@echo ""
	@echo "Running test application..."
	@./vf_test

clean:
	@echo "Cleaning up..."
	rm -f vf_test test.db

test:
	@echo "Running tests..."
	go test -v

help:
	@echo "Available targets:"
	@echo "  make build  - Build the test application"
	@echo "  make run    - Build and run the test application"
	@echo "  make clean  - Remove built files and test database"
	@echo "  make test   - Run tests"
	@echo "  make all    - Build and run (default)"
