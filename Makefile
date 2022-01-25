run:
	echo "Running application"
	swag init --generalInfo cmd/main.go --output api/openapi
	go run cmd/main.go