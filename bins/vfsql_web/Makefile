# filepath: /home/lxk/Desktop/cirm/bins/vfsql_web/Makefile
.PHONY: all build run clean help setup

all: setup build run

setup:
	@if [ ! -e "web" ] && [ -d "../vfsql/web" ]; then \
		echo "Creating symlink to web directory..."; \
		ln -sf ../vfsql/web web; \
	fi

build:
	@echo "Building VFSQL Web Server..."
	go build -o vfsql_web .

run: setup build
	@echo ""
	@echo "Starting server..."
	@./vfsql_web

dev: setup build
	@./vfsql_web -db dev.db

clean:
	@echo "Cleaning up..."
	rm -f vfsql_web vfsql.db dev.db
	@if [ -L "web" ]; then rm -f web; fi

help:
	@echo "Available targets:"
	@echo "  make setup  - Create web directory symlink"
	@echo "  make build  - Build the web server"
	@echo "  make run    - Build and run the server"
	@echo "  make dev    - Run with dev.db database"
	@echo "  make clean  - Remove built files and databases"
	@echo "  make all    - Setup, build and run (default)"
