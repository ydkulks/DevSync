check:
	@echo "Checking if all the agent providers have proper metadata format..."
	cd ./scripts/ && go run ./check-format.go
	@echo "Checking if all the agent providers are in sync..."
	cd ./scripts/ && go run ./check-sync.go
