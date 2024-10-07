air:
	@echo "Running the Go project using Air..."
	air

run:
	@echo "Running the Go project..."
	go run .

ngrok:
	@echo "Running the Go project using Ngrok..."
	ngrok http 8080