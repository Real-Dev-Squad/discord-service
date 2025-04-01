air:
	@echo "Running the Go project using Air..."
	air

run:
	@echo "Running the Go project..."
	go run .

ngrok:
	@echo "Running the Go project using Ngrok..."
	ngrok http 8080

tidy:
	@echo "Running the Go project tidy..."
	@go mod tidy

download:
	@echo "Running the Go project tidy..."
	@go mod download

fmt:
	@echo "Running the Go project fmt..."
	@go fmt ./...

test:
	@echo "Running the Go project tests..."
	@go list ./... | grep -v "/config$$" | grep -v "/routes$$" | xargs go test -v -coverprofile=coverage.out

coverage:
	@echo "Running the Go project tests with coverage..."
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

clean:
	@echo "Cleaning the Go project..."
	@rm -rf coverage
	@rm -rf coverage.out
	@rm -rf coverage.html

test-cover:
ifeq ($(FORCE),1)
	@echo "Force flag detected. Cleaning before tests..."
	@$(MAKE) clean
endif
	@$(MAKE) test
	@$(MAKE) coverage
	@echo "Tests completed and coverage report generated."


	
