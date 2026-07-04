
# Build the coroki application into the bin/ directory
# With the name `coroki`
build:
	go build -o bin/coroki cmd/coroki/main.go

# Runs the coroki application.
run:
	go run cmd/coroki/main.go
